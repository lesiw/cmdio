package sys

import (
	"bytes"
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

func Command(args ...string) io.ReadWriter {
	return defaultBox.Command(args...)
}

type cmd struct {
	ctx  context.Context
	cmd  *exec.Cmd
	log  []byte
	env  map[string]string
	code int

	start func() error
	wait  func() error

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
	c.start = sync.OnceValue(c._start)
	c.wait = sync.OnceValue(c._wait)
}

func (c *cmd) _start() (err error) {
	if c.cmd.Stdin == nil {
		if c.writer, err = c.cmd.StdinPipe(); err != nil {
			return fmt.Errorf("failed to pipe stdin: %w", err)
		}
	}
	if c.cmd.Stdout == nil {
		if c.reader, err = c.cmd.StdoutPipe(); err != nil {
			return fmt.Errorf("failed to pipe stdout: %w", err)
		}
	}
	if c.cmd.Stderr == nil {
		if c.logger, err = c.cmd.StderrPipe(); err != nil {
			return fmt.Errorf("failed to pipe stderr: %w", err)
		}
	}
	return cmdio.NewError(c.cmd.Start(), c)
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
	return bytes.NewReader(c.log)
}

func (c *cmd) Code() int {
	return c.code
}

func (c *cmd) _wait() error {
	if c.logger != nil {
		buf, err := io.ReadAll(c.logger)
		if err != nil {
			return fmt.Errorf("failed to read stderr: %w", err)
		}
		c.log = buf
	}
	err := c.cmd.Wait() // Closes pipes.
	if err == nil {
		return nil
	}
	ee := new(exec.ExitError)
	if errors.As(err, &ee) {
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
