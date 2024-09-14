// Package cmdio provides an interface to treat commands as an I/O stream.
package cmdio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

// Trace is an [io.Writer] to which stream tracing information is written.
// To disable, set this variable to [io.Discard].
var Trace io.Writer = newPrefixWriter("+ ", Stderr)

// A Logger has a readable log.
type Logger interface {
	Log() io.Reader
}

// A Coder has an exit code.
type Coder interface {
	Code() int
}

// An Attacher can be connected directly to the controlling terminal.
// An attached stream cannot be written to.
// It must be readable exactly once, during which time the read must block for
// the duration of the command, then exit with 0 bytes read.
type Attacher interface {
	Attach() error
}

// Run attaches a command to the controlling terminal and executes it.
func Run(cmd io.Reader) error {
	fmt.Fprintln(Trace, cmd)
	a, ok := cmd.(Attacher)
	if !ok {
		// If this command does not implement Attacher, stream it to stdout
		// (and stderr, if applicable).
		if l, ok := cmd.(Logger); ok {
			go io.Copy(Stderr, l.Log())
		}
		_, err := io.Copy(Stdout, cmd)
		return NewError(err, readWriter(cmd))
	}
	if err := a.Attach(); err != nil {
		return err
	}
	_, err := cmd.Read(nil)
	if err == io.EOF {
		err = nil
	}
	return NewError(err, readWriter(cmd))
}

// Get executes a command and captures its output.
// Result is never nil, even if error is not nil.
// Checking Result.Code > 0 is not sufficient, as it will default to 0 even in
// cases where the command does not successfully complete, such as a "command
// not found" error.
func Get(cmd io.Reader) (*Result, error) {
	fmt.Fprintln(Trace, cmd)

	r := new(Result)
	var errs []error

	buf, err := io.ReadAll(cmd)
	if err != nil {
		errs = append(errs, err)
	}
	r.Cmd = readWriter(cmd)
	r.Out = strings.Trim(string(buf), "\n")
	if l, ok := cmd.(Logger); ok {
		logbuf, err := io.ReadAll(l.Log())
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to read cmd log: %w", err))
		}
		r.Log = strings.Trim(string(logbuf), "\n")
	}
	if c, ok := cmd.(Coder); ok {
		r.Code = c.Code()
	}
	return r, errors.Join(errs...)
}
