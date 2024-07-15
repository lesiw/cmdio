package cmd

import (
	"context"
	"io"
	"os/exec"

	"lesiw.io/cmdio"
)

var defaultBox = cmdio.NewBox(&cmdBox{})

type cmdBox struct {
	env map[string]string
}

func (b *cmdBox) Command(
	ctx context.Context, args ...string,
) io.ReadWriter {
	s := &cmd{
		ctx: ctx,
		cmd: exec.CommandContext(ctx, args[0], args[1:]...),
		env: b.env,
	}
	s.init()
	return s
}

func (b *cmdBox) Env(k string) string {
	if b.env == nil {
		return ""
	}
	return b.env[k]
}

func (b *cmdBox) Setenv(k, v string) {
	if b.env == nil {
		b.env = make(map[string]string)
	}
	b.env[k] = v
}

func Env(env map[string]string) *cmdio.Box {
	n := new(cmdBox)
	for k, v := range env {
		n.Setenv(k, v)
	}
	return cmdio.NewBox(n)
}

func Run(args ...string) error {
	return defaultBox.Run(args...)
}

func MustRun(args ...string) {
	defaultBox.MustRun(args...)
}

func Get(args ...string) (*cmdio.CmdResult, error) {
	return defaultBox.Get(args...)
}

func MustGet(args ...string) *cmdio.CmdResult {
	return defaultBox.MustGet(args...)
}
