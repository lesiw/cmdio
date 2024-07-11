package cmd

import (
	"context"
	"io"
	"os/exec"

	"lesiw.io/cmdio"
)

var defaultNS = cmdio.NewNS(&cmdNS{})

type cmdNS struct {
	env map[string]string
}

func (ns *cmdNS) CommandContext(
	ctx context.Context, args ...string,
) io.ReadWriter {
	s := &cmd{
		ctx: ctx,
		cmd: exec.CommandContext(ctx, args[0], args[1:]...),
		env: ns.env,
	}
	s.init()
	return s
}

func (ns *cmdNS) Env(k string) string {
	if ns.env == nil {
		return ""
	}
	return ns.env[k]
}

func (ns *cmdNS) Setenv(k, v string) {
	if ns.env == nil {
		ns.env = make(map[string]string)
	}
	ns.env[k] = v
}

func Env(env map[string]string) *cmdio.NS {
	n := new(cmdNS)
	for k, v := range env {
		n.Setenv(k, v)
	}
	return cmdio.NewNS(n)
}

func Run(args ...string) error {
	return defaultNS.Run(args...)
}

func MustRun(args ...string) {
	defaultNS.MustRun(args...)
}

func Check(args ...string) (*cmdio.CmdResult, error) {
	return defaultNS.Check(args...)
}

func MustCheck(args ...string) *cmdio.CmdResult {
	return defaultNS.MustCheck(args...)
}

func Get(args ...string) (*cmdio.CmdResult, error) {
	return defaultNS.Get(args...)
}

func MustGet(args ...string) *cmdio.CmdResult {
	return defaultNS.MustGet(args...)
}
