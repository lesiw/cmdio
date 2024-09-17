package cmdio

import (
	"bufio"
	"context"
	"io"
	"strings"
)

// A Commander is an interface representing a computer.
type Commander interface {
	Command(context.Context, map[string]string, ...string) io.ReadWriter
}

// An Enver retrieves environment variables.
type Enver interface {
	Env(string) string
}

// Box is an abstraction of a computer.
// It represents any entity that can execute commands.
// This might be the system, a container, a remote server, or something else.
type Box struct {
	ctx context.Context
	env map[string]string
	cdr Commander
}

// NewBox creates a new box.
func NewBox(ctx context.Context, env map[string]string, cmd Commander) *Box {
	return &Box{ctx, env, cmd}
}

// WithEnv creates a new box with the provided env.
// The new box will share the same context and commander as its parent.
func (b *Box) WithEnv(env map[string]string) *Box {
	return &Box{b.ctx, env, b.cdr}
}

// WithContext creates a new box with the provided [context.Context].
// The new box will share the same environment and commander as its parent.
func (b *Box) WithContext(ctx context.Context) *Box {
	return &Box{ctx, b.env, b.cdr}
}

// Command returns an command as an io.ReadWriter.
// The command will not be executed until it is read.
func (b *Box) Command(args ...string) io.ReadWriter {
	ctx := b.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return b.cdr.Command(ctx, b.env, args...)
}

// Run runs a command.
func (b *Box) Run(args ...string) error {
	return Run(b.Command(args...))
}

// MustRun runs a command and panics on failure.
func (b *Box) MustRun(args ...string) {
	must(b.Run(args...))
}

// Get runs a command and captures its output in a [Result].
func (b *Box) Get(args ...string) (*Result, error) {
	return Get(b.Command(args...))
}

// MustGet runs a command and captures its output in a [Result].
// It panics if the command fails.
func (b *Box) MustGet(args ...string) *Result {
	return must1(b.Get(args...))
}

// Close closes the underlying Commander if it is an [io.Closer].
func (b *Box) Close() error {
	if closer, ok := b.cdr.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Env returns the value of an environment variable.
func (b *Box) Env(name string) (value string) {
	if enver, ok := b.cdr.(Enver); ok {
		return enver.Env(name)
	}
	scanner := bufio.NewScanner(b.Command("env"))
	for scanner.Scan() {
		line := scanner.Text()
		k, v, ok := strings.Cut(line, "=")
		if ok && k == name {
			return v
		}
	}
	return
}
