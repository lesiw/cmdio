package ctr

import (
	"context"

	"lesiw.io/cmdio"
)

type cmd struct {
	cmdio.Command
	cdr *cdr
	ctx context.Context
	env map[string]string
	arg []string
}

func newCmd(
	cdr *cdr, ctx context.Context, env map[string]string, args ...string,
) cmdio.Command {
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
	if a, ok := c.Command.(cmdio.Attacher); ok {
		return a.Attach()
	}
	return nil
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
	c.Command = c.cdr.rnr.Commander.Command(c.ctx, nil, cmd...)
}
