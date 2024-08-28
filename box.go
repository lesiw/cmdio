package cmdio

import (
	"context"
	"io"
)

type Commander interface {
	Command(context.Context, ...string) io.ReadWriter
}

// Box represents any entity that can execute commands.
// This might be the system, a container, a remote server, or something else.
type Box struct {
	Commander
	ctx context.Context
}

func NewBox(c Commander) *Box {
	return &Box{c, nil}
}

func NewBoxContext(c Commander, ctx context.Context) *Box {
	return &Box{c, ctx}
}

func (b *Box) Command(args ...string) io.ReadWriter {
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return b.Commander.Command(ctx, args...)
}

// Run runs a command.
func (b *Box) Run(args ...string) error {
	return Run(b.Command(args...))
}

// MustRun runs a command and panics on failure.
func (b *Box) MustRun(args ...string) {
	must(b.Run(args...))
}

// Get runs a command and captures its output in a [CmdResult].
func (b *Box) Get(args ...string) (*CmdResult, error) {
	return Get(b.Command(args...))
}

// MustGet runs a command and captures its output in a [CmdResult].
// It panics if the command fails.
func (b *Box) MustGet(args ...string) *CmdResult {
	return must1(b.Get(args...))
}

// WithContext returns a new copy of the provided Box with a context.
func WithContext(b *Box, ctx context.Context) *Box {
	return &Box{b.Commander, ctx}
}
