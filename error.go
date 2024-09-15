package cmdio

import (
	"bytes"
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
	if l, ok := e.err.(Logger); ok {
		buf, err := io.ReadAll(l.Log())
		if err == nil && len(buf) > 0 {
			fmt.Fprintf(w, "\nlog:\n---\n%s\n---\n", string(buf))
		}
	}
}

func (e *Error) Unwrap() error {
	return e.err
}

type logError struct {
	err error
	log []byte
}

func (e *logError) Error() string {
	return e.err.Error()
}

func (e *logError) Log() io.Reader {
	return bytes.NewReader(e.log)
}

func (e *logError) Unwrap() error {
	return e.err
}
