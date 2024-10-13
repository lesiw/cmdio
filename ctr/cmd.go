package ctr

import (
	"context"
	"fmt"
	"io"

	"lesiw.io/cmdio"
)

type cmd struct {
	io.ReadWriter
	cdr *cdr
	ctx context.Context
	env map[string]string
	arg []string
}

func newCmd(
	cdr *cdr, ctx context.Context, env map[string]string, args ...string,
) io.ReadWriter {
	c := &cmd{
		ctx: ctx,
		env: env,
		cdr: cdr,
		arg: args,
	}
	c.setCmd(false)
	return c
}

func (c *cmd) Attach() error {
	c.setCmd(true)
	if a, ok := c.ReadWriter.(cmdio.Attacher); ok {
		return a.Attach()
	}
	return nil
}

func (c *cmd) Close() error {
	if cl, ok := c.ReadWriter.(io.Closer); ok {
		return cl.Close()
	}
	return nil
}

func (c *cmd) String() string {
	if s, ok := c.ReadWriter.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("<%T>", c)
}

func (c *cmd) setCmd(attach bool) {
	cmd := []string{"container", "exec", "-i"}
	if attach {
		cmd = append(cmd, "-t")
	}
	for k, v := range c.env {
		if k == "PWD" {
			cmd = append(cmd, "-w", c.env["PWD"])
		} else {
			cmd = append(cmd, "-e", k+"="+v)
		}
	}
	cmd = append(cmd, c.cdr.ctrid)
	cmd = append(cmd, c.arg...)
	c.ReadWriter = c.cdr.rnr.Commander.Command(c.ctx, nil, cmd...)
}
