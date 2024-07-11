package cmdio

import "io"

// Error describes an error produced by the cmdio package that can be recovered
// with [Recover].
type Error struct {
	err  error
	cmd  io.ReadWriter
	Code int
	Log  string
}

func NewError(err error, cmd io.ReadWriter) error {
	if err == nil {
		return nil
	}
	return &Error{err: err, cmd: cmd}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}
