// Package cmdio provides an interface to treat commands as an I/O stream.
package cmdio

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Trace is an [io.Writer] to which stream tracing information is written.
// To disable, set this variable to [io.Discard].
var Trace io.Writer = os.Stderr

// AttachReadWriter is an I/O abstraction over a system command.
type AttachReadWriter interface {
	io.ReadWriter

	// Attach connects the stream to the controlling terminal.
	// An attached stream cannot be written to.
	// It can be read from exactly once, during which time the read will block
	// for the duration of the command, then exit, having read 0 bytes.
	Attach() error
}

// Run attaches a stream to the controlling terminal and executes it.
func Run(s AttachReadWriter) error {
	fmt.Fprintln(Trace, "+", s)
	if err := s.Attach(); err != nil {
		return err
	}
	if _, err := s.Read(nil); err != nil {
		return NewError(err, s)
	}
	return nil
}

// Check executes a stream, capturing its output.
// If the command runs and exits non-zero, it will not return an error.
// If the command fails to run, however, it will return an error.
func Check(s io.ReadWriter) (*CmdResult, error) {
	fmt.Fprintln(Trace, "+", s)
	buf, err := io.ReadAll(s)
	r := &CmdResult{
		Cmd:    s,
		Output: strings.Trim(string(buf), "\n"),
	}
	if err != nil {
		se := new(Error)
		if errors.As(err, &se) {
			if se.Code > 0 {
				r.Code = se.Code
				return r, nil
			}
			return nil, se
		} else {
			return nil, NewError(err, s)
		}
	}
	r.Ok = true
	return r, nil
}

// Get executes a stream and captures its output.
func Get(s io.ReadWriter) (*CmdResult, error) {
	fmt.Fprintln(Trace, "+", s)
	buf, err := io.ReadAll(s)
	if err != nil {
		return nil, err
	}
	r := &CmdResult{
		Ok:     true,
		Cmd:    s,
		Output: strings.Trim(string(buf), "\n"),
	}
	return r, nil
}
