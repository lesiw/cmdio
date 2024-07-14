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
}

func NewBox(c Commander) *Box {
	return &Box{c}
}

func (b *Box) Command(args ...string) io.ReadWriter {
	return b.CommandContext(context.Background(), args...)
}

func (b *Box) CommandContext(
	ctx context.Context, args ...string,
) io.ReadWriter {
	return b.Commander.Command(context.Background(), args...)
}

// Run runs a command.
func (b *Box) Run(args ...string) error {
	return Run(b.Command(args...).(AttachReadWriter))
}

// MustRun runs a command and panics on failure.
func (b *Box) MustRun(args ...string) {
	must(b.Run(args...))
}

// Check runs a command and returns a [CmdResult].
// It will not return an error if the command exits with a non-zero exit status.
// However, it will return an error if it fails to run.
func (b *Box) Check(args ...string) (*CmdResult, error) {
	return Check(b.Command(args...))
}

// Check runs a command and returns a [CmdResult].
// It will not panic if the command exits with a non-zero exit status.
// However, it will panic if it fails to run.
func (b *Box) MustCheck(args ...string) *CmdResult {
	return must1(b.Check(args...))
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
