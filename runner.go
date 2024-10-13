package cmdio

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
)

// A Runner runs commands.
type Runner struct {
	ctx context.Context
	env map[string]string
	Commander
}

// NewRunner instantiates a new [Runner].
func NewRunner(
	ctx context.Context, env map[string]string, cdr Commander,
) *Runner {
	return &Runner{ctx, env, cdr}
}

// WithEnv creates a new runner with the provided env.
// The new runner will share the same context and commander as its parent.
//
// PWD conventionally sets the working directory.
func (rnr *Runner) WithEnv(env map[string]string) *Runner {
	env2 := make(map[string]string, len(rnr.env))
	for k, v := range rnr.env {
		env2[k] = v
	}
	for k, v := range env {
		env2[k] = v
	}
	return &Runner{rnr.ctx, env2, rnr.Commander}
}

// WithContext creates a new runner with the provided [context.Context].
// The new runner will share the same environment and commander as its parent.
func (rnr *Runner) WithContext(ctx context.Context) *Runner {
	env := make(map[string]string, len(rnr.env))
	for k, v := range rnr.env {
		env[k] = v
	}
	return &Runner{ctx, env, rnr.Commander}
}

// Command instantiates a command as an [io.ReadWriter].
//
// The command will not be executed until the first time it is read or written
// to.
func (rnr *Runner) Command(args ...string) io.ReadWriter {
	ctx := rnr.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	return rnr.Commander.Command(ctx, rnr.env, args...)
}

// Run attaches a command to the controlling terminal and executes it.
func (rnr *Runner) Run(args ...string) error {
	return run(rnr.Command(args...))
}

// MustRun runs a command and panics on failure.
func (rnr *Runner) MustRun(args ...string) {
	must(rnr.Run(args...))
}

// Get executes a command and captures the output in a [Result].
//
// Note that checking Result.Code > 0 is not sufficient to determine that the
// command executed successfully. Commands may choose not to implement [Coder],
// and commands that fail to execute because they cannot be found will have no
// exit code.
func (rnr *Runner) Get(args ...string) (Result, error) {
	return get(rnr.Command(args...))
}

// MustGet runs a command and captures its output in a [Result].
// It panics with diagnostic output if the command fails.
func (rnr *Runner) MustGet(args ...string) Result {
	r, err := rnr.Get(args...)
	if err != nil {
		panic(err.Error() + "\n" +
			"out:" + fmtout(r.Out) +
			"log:" + fmtout(r.Log) +
			"code: " + strconv.Itoa(r.Code))
	}
	return r
}

// Close closes the underlying [Commander] if it is an [io.Closer].
func (rnr *Runner) Close() error {
	if closer, ok := rnr.Commander.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// Env returns the value of an environment variable.
//
// By default, it parses the output of an env command. [Commander]
// implementations may customize this behavior by implementing [Enver].
func (rnr *Runner) Env(name string) (value string) {
	if enver, ok := rnr.Commander.(Enver); ok {
		return enver.Env(name)
	}
	scanner := bufio.NewScanner(rnr.Command("env"))
	for scanner.Scan() {
		line := scanner.Text()
		k, v, ok := strings.Cut(line, "=")
		if ok && k == name {
			return v
		}
	}
	return
}
