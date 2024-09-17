package sys

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	"lesiw.io/cmdio"
)

type cmd struct {
	ctx  context.Context
	cmd  *exec.Cmd
	env  map[string]string
	code int

	cmdwait chan error

	start func() error
	wait  func() error

	writer io.WriteCloser
	reader io.ReadCloser
	logger io.ReadCloser

	ready chan struct{}

	closers []io.Closer
}

func (c *cmd) Attach() error {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
	return nil
}

func newCmd(ctx context.Context, env map[string]string, args ...string) *cmd {
	c := &cmd{}
	c.ctx = ctx
	c.cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	c.env = env
	if dir, ok := env["PWD"]; ok {
		c.cmd.Dir = dir
	}
	c.cmd.Env = os.Environ()
	for k, v := range env {
		c.cmd.Env = append(c.cmd.Env, k+"="+v)
	}
	c.start = sync.OnceValue(c.startFunc)
	c.wait = sync.OnceValue(c.waitFunc)
	c.cmdwait = make(chan error)
	c.ready = make(chan struct{})
	return c
}

func (c *cmd) startFunc() error {
	if c.cmd.Stdin == nil {
		w, err := c.cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("failed to pipe stdin: %w", err)
		}
		c.writer = w
	}
	if c.cmd.Stdout == nil {
		r, w := io.Pipe()
		c.reader = r
		c.cmd.Stdout = w
		c.closers = append(c.closers, w)
	}
	if c.cmd.Stderr == nil {
		r, w := io.Pipe()
		c.logger = r
		c.cmd.Stderr = w
		c.closers = append(c.closers, w)
	}
	close(c.ready)
	if err := c.cmd.Start(); err != nil {
		for _, cl := range c.closers {
			_ = cl.Close() // Best effort.
		}
		return cmdio.NewError(err, c)
	}
	go func() {
		err := c.cmd.Wait()
		for _, cl := range c.closers {
			if err1 := cl.Close(); err == nil {
				err = err1
			}
		}
		c.cmdwait <- err
	}()
	return nil
}

func (c *cmd) Write(bytes []byte) (int, error) {
	if err := c.start(); err != nil {
		return 0, err
	}
	if c.writer == nil {
		return 0, nil
	}
	n, err := c.writer.Write(bytes)
	if err != nil {
		return n, cmdio.NewError(fmt.Errorf(`failed write: %w`, err), c)
	}
	return n, nil
}

func (c *cmd) Close() error {
	if c.writer == nil {
		return nil
	}
	if err := c.writer.Close(); err != nil {
		return cmdio.NewError(fmt.Errorf(`failed close: %w`, err), c)
	}
	return nil
}

func (c *cmd) Read(bytes []byte) (int, error) {
	if err := c.start(); err != nil {
		return 0, err
	}
	ch := make(chan ioret)
	var n int
	var err error
	if c.reader == nil {
		err = io.EOF
		goto skipread
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

skipread:
	if err != nil {
		if err1 := c.wait(); err1 != nil {
			err = err1
		}
	}
	return n, err
}

func (c *cmd) Log() io.Reader {
	<-c.ready
	return c.logger
}

func (c *cmd) Code() int {
	return c.code
}

func (c *cmd) waitFunc() error {
	err := <-c.cmdwait
	if ee := new(exec.ExitError); errors.As(err, &ee) {
		c.code = ee.ExitCode()
	}
	return cmdio.NewError(err, c)
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
