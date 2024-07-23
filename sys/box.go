package sys

import (
	"context"
	"io"
	"os/exec"

	"lesiw.io/cmdio"
)

var defaultBox = cmdio.NewBox(&sysbox{})

type sysbox struct {
	env map[string]string
}

type Enver interface {
	Env(string) string
	Setenv(string, string)
}

func (b *sysbox) Command(ctx context.Context, args ...string) io.ReadWriter {
	s := &cmd{
		ctx: ctx,
		cmd: exec.CommandContext(ctx, args[0], args[1:]...),
		env: b.env,
	}
	s.init()
	return s
}

func (b *sysbox) Env(k string) string {
	if b.env == nil {
		return ""
	}
	return b.env[k]
}

func (b *sysbox) Setenv(k, v string) {
	if b.env == nil {
		b.env = make(map[string]string)
	}
	b.env[k] = v
}

func Env(env map[string]string) *cmdio.Box {
	return WithEnv(new(sysbox), env)
}

func WithEnv(b Enver, env map[string]string) *cmdio.Box {
	for k, v := range env {
		b.Setenv(k, v)
	}
	return cmdio.NewBox(b.(cmdio.Commander))
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
