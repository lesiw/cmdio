package cmdio

import (
	"fmt"
	"io"
)

// Error describes an error produced by the cmdio package that can be recovered
// with [Recover].
type Error struct {
	err error
	cmd io.ReadWriter
}

func NewError(err error, cmd io.ReadWriter) error {
	if err == nil {
		return nil
	}
	return &Error{err, cmd}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Print(w io.Writer) {
	fmt.Fprintf(w, "exec failed: %v: %s\n", e.cmd, e.Error())
	if l, ok := e.cmd.(Logger); ok {
		fmt.Fprintf(w, "\nlog:\n---\n")
		if _, err := io.Copy(w, l.Log()); err != nil {
			fmt.Fprintf(w, "--- error reading log: %s ---", err)
		}
		fmt.Fprintf(w, "\n---\n")
	}
}

func (e *Error) Unwrap() error {
	return e.err
}
