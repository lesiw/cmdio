// Package cmdio provides portable interfaces for commands and command runners.
//
// A command is an [io.ReadWriter]. Writing to a command writes to its standard
// input. Reading from a command reads from its standard output. Commands may
// optionally implement [Logger] to capture standard error and [Coder] to
// represent exit codes.
//
// Commands are instantiated by a [Runner]. This package contains several
// Runner implementations: [lesiw.io/cmdio/sys], which runs commands on the
// local system; [lesiw.io/cmdio/ctr], which runs commands in containers; and
// [lesiw.io/cmdio/sub], which runs commands as subcommands.
//
// While most of this package is written to support traditional Go error
// handling, Must-type functions, such as [Runner.MustRun] and [MustPipe], are
// provided to support a script-like programming style, where failures result
// in panics.
package cmdio

import (
	"context"
	"fmt"
	"io"
	"os"

	"lesiw.io/prefix"
)

// A Commander instantiates commands.
//
// The Command function accepts a [context.Context], a map of environment
// variables, and a variable number of arguments representing the command
// itself. It returns a [Command].
type Commander interface {
	Command(
		ctx context.Context,
		env map[string]string,
		arg ...string,
	) (cmd Command)
}

// An Enver has environment variables.
//
// A [Commander] that also implements this interface will call Env to retrieve
// environment variables.
type Enver interface {
	Env(name string) (value string)
}

// A Logger accepts an [io.Writer] for logging diagnostic information.
//
// Implementing this interface is the idiomatic way for commands to represent
// standard error.
type Logger interface {
	Log(io.Writer)
}

// A Coder has an exit code.
//
// Implementing this interface is the idiomatic way for commands to represent
// exit codes.
type Coder interface {
	Code() int
}

// An Attacher can be connected directly to the controlling terminal.
// An attached command cannot be written to.
// It must be readable exactly once. The read must block for the duration of
// command execution, after which it must exit with 0 bytes read.
type Attacher interface {
	Attach() error
}

// A [Command] is the broadest possible command interface.
//
// Commands must not begin execution until the first time they are read from or
// written to. They must return [io.EOF] once execution has completed and all
// output has been consumed.
//
// In general, the Write method will correspond to standard in, the Read
// method will correspond to standard out, and an [io.Writer] may be passed
// to Log for handling standard error.
type Command interface {
	io.ReadWriteCloser
	fmt.Stringer
	Attacher
	Coder
	Logger
}

// Trace is an [io.Writer] to which command tracing information is written.
// To disable tracing, set this variable to [io.Discard].
var Trace io.Writer = prefix.NewWriter("+ ", stderr)

var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)
