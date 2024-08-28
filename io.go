// Package cmdio provides an interface to treat commands as an I/O stream.
package cmdio

import (
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

// An attacher can be connected directly to the controlling terminal.
// An attached stream cannot be written to.
// It must be readable exactly once, during which time the read must block for
// the duration of the command, then exit with 0 bytes read.
type Attacher interface {
	Attach() error
}

// Run attaches a command to the controlling terminal and executes it.
func Run(cmd io.Reader) error {
	fmt.Fprintln(Trace, "+", cmd)
	a, ok := cmd.(Attacher)
	if !ok {
		// If this command does not implement Attacher, stream it to stdout.
		_, err := io.Copy(Stdout, cmd)
		return NewError(err, cmd)
	}
	if err := a.Attach(); err != nil {
		return err
	}
	_, err := cmd.Read(nil)
	return NewError(err, cmd)
}

// Get executes a command and captures its output.
func Get(cmd io.Reader) (*CmdResult, error) {
	fmt.Fprintln(Trace, "+", cmd)
	buf, err := io.ReadAll(cmd)
	if err != nil {
		return nil, err
	}
	r := &CmdResult{
		Cmd:    cmd,
		Output: strings.Trim(string(buf), "\n"),
	}
	return r, nil
}
