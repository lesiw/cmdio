package cmd

import (
	"cmp"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"lesiw.io/cmdio"
)

func Command(args ...string) io.ReadWriter {
	return defaultBox.Command(args...)
}

func CommandContext(ctx context.Context, args ...string) io.ReadWriter {
	return defaultBox.CommandContext(ctx, args...)
}

type cmd struct {
	ctx context.Context
	cmd *exec.Cmd
	log []byte
	env map[string]string

	starter sync.Once
	waiter  sync.Once

	writer io.WriteCloser
	reader io.ReadCloser
	logger io.ReadCloser
}

func (c *cmd) Attach() error {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return nil
}

func (c *cmd) init() {
	c.cmd.Env = os.Environ()
	for k, v := range c.env {
		c.cmd.Env = append(c.cmd.Env, k+"="+v)
	}
}

func (c *cmd) start() {
	if c.cmd.Stdin == nil {
		c.writer = must1(c.cmd.StdinPipe())
	}
	if c.cmd.Stdout == nil {
		c.reader = must1(c.cmd.StdoutPipe())
	}
	if c.cmd.Stderr == nil {
		c.logger = must1(c.cmd.StderrPipe())
	}
	must(cmdio.NewError(c.cmd.Start(), c))
}

func (c *cmd) Write(bytes []byte) (int, error) {
	c.starter.Do(c.start)
	if c.writer == nil {
		return 0, nil
	}
	return c.writer.Write(bytes)
}

func (c *cmd) Close() error {
	if c.writer == nil {
		return nil
	}
	return c.writer.Close()
}

func (c *cmd) Read(bytes []byte) (int, error) {
	c.starter.Do(c.start)
	ch := make(chan ioret)
	var n int
	var err error
	if c.reader == nil {
		goto nilreader
	}

	go func() {
		n, err := c.reader.Read(bytes)
		ch <- ioret{n, err}
	}()
	select {
	case <-c.ctx.Done():
		n = 0
		err = io.EOF
	case ret := <-ch:
		n = ret.n
		err = ret.err
	}

nilreader:
	if err != nil || n == 0 {
		c.waiter.Do(func() {
			if err1 := c.wait(); err1 != nil {
				err = err1
			}
		})
	}
	return n, err
}

func (c *cmd) wait() error {
	var log []byte
	if c.logger != nil {
		log = must1(io.ReadAll(c.logger))
	}
	err := c.cmd.Wait() // Closes pipes.
	if err == nil {
		return nil
	}
	ee := new(exec.ExitError)
	ce := cmdio.NewError(err, c).(*cmdio.Error)
	if errors.As(err, &ee) {
		ce.Code = ee.ExitCode()
	}
	ce.Log = strings.TrimRight(string(log), "\n")
	return ce
}

func (c *cmd) String() string {
	ret := new(strings.Builder)
	for _, k := range sortkeys(c.env) {
		ret.WriteString(k + "=" + c.env[k] + " ")
	}
	ret.WriteString(shJoin(c.cmd.Args))
	return ret.String()
}

type ioret struct {
	n   int
	err error
}

func sortkeys[K cmp.Ordered, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	var i int
	for k := range m {
		keys[i] = k
		i++
	}
	slices.Sort(keys)
	return keys
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must1[T any](r0 T, err error) T {
	if err != nil {
		panic(err)
	}
	return r0
}
