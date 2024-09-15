// Package cmdio provides an interface to treat commands as an I/O stream.
package cmdio

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"
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
	var wg errgroup.Group
	out := make(chan string)
	log := make(chan string)

	if l, ok := cmd.(Logger); ok {
		wg.Go(func() error {
			buf, err := io.ReadAll(l.Log())
			log <- strings.Trim(string(buf), "\n")
			if err != nil {
				return fmt.Errorf("failed to read cmd log: %w", err)
			}
			return nil
		})
	} else {
		close(out)
	}
	wg.Go(func() error {
		buf, err := io.ReadAll(cmd)
		out <- strings.Trim(string(buf), "\n")
		return err
	})

	r.Cmd = readWriter(cmd)
	r.Out = <-out
	r.Log = <-log
	err := wg.Wait()
	if c, ok := cmd.(Coder); ok {
		r.Code = c.Code()
	}

	if err != nil && r.Log != "" {
		err = NewError(&logError{err, []byte(r.Log)}, r.Cmd)
	}

	return r, err
}
