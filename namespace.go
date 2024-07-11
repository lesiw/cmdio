package cmdio

import (
	"context"
	"io"
)

type Namespace interface {
	CommandContext(context.Context, ...string) io.ReadWriter
}

type NS struct {
	Namespace
}

func NewNS(ns Namespace) *NS {
	return &NS{ns}
}

func (ns *NS) Command(args ...string) io.ReadWriter {
	return ns.CommandContext(context.Background(), args...)
}

// Run runs a command.
func (ns *NS) Run(args ...string) error {
	return Run(ns.Command(args...).(AttachReadWriter))
}

// MustRun runs a command and panics on failure.
func (ns *NS) MustRun(args ...string) {
	must(ns.Run(args...))
}

// Check runs a command and returns a [CmdResult].
// It will not return an error if the command exits with a non-zero exit status.
// However, it will return an error if it fails to run.
func (ns *NS) Check(args ...string) (*CmdResult, error) {
	return Check(ns.Command(args...))
}

// Check runs a command and returns a [CmdResult].
// It will not panic if the command exits with a non-zero exit status.
// However, it will panic if it fails to run.
func (ns *NS) MustCheck(args ...string) *CmdResult {
	return must1(ns.Check(args...))
}

// Get runs a command and captures its output in a [CmdResult].
func (ns *NS) Get(args ...string) (*CmdResult, error) {
	return Get(ns.Command(args...))
}

// MustGet runs a command and captures its output in a [CmdResult].
// It panics if the command fails.
func (ns *NS) MustGet(args ...string) *CmdResult {
	return must1(ns.Get(args...))
}
