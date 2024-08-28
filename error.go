package cmdio

import (
	"fmt"
	"io"
)

// Error describes an error produced by the cmdio package that can be recovered
// with [Recover].
type Error struct {
	err  error
	Cmd  any
	Code int
	Log  string
}

func NewError(err error, cmd any) error {
	if err == nil {
		return nil
	}
	return &Error{err: err, Cmd: cmd}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Print(w io.Writer) {
	fmt.Fprintf(w, "exec failed: %v: %s\n", e.Cmd, e.Error())
	if e.Log != "" {
		fmt.Fprintf(w, "\nstderr:\n---\n%s\n---\n", e.Log)
	}
}

func (e *Error) Unwrap() error {
	return e.err
}
