package cmdio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// RunPipe chains commands together.
func RunPipe(src io.Reader, cmds ...io.ReadWriter) error {
	all := make([]any, len(cmds)+1)
	all[0] = src
	for i, s := range cmds {
		all[i+1] = s
	}
	printCmds(all...)
	_, err := Copy(nopCloser{os.Stdout}, src, cmds...)
	return err
}

// MustRunPipe chains streams together and panics on failure.
func MustRunPipe(src io.Reader, cmds ...io.ReadWriter) {
	must(RunPipe(src, cmds...))
}

// CheckPipe chains streams together.
// On failure, the failing stream's information will be returned.
func CheckPipe(src io.Reader, cmds ...io.ReadWriter) (*CmdResult, error) {
	all := make([]any, len(cmds)+1)
	all[0] = src
	for i, s := range cmds {
		all[i+1] = s
	}
	printCmds(all...)
	r := new(CmdResult)
	if _, err := Copy(io.Discard, src, cmds...); err != nil {
		if ce := new(Error); errors.As(err, &ce) {
			if ce.Code > 0 {
				r.Cmd = ce.cmd
				r.Code = ce.Code
				return r, nil
			}
			return nil, ce
		} else if ce := new(CopyError); errors.As(err, &ce) {
			if s, ok := ce.Reader.(io.ReadWriter); ok {
				return nil, NewError(err, s)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	r.Ok = true
	return r, nil
}

// MustCheckPipe chains commands together, panicking if a command cannot start.
// On failure, the failing stream's information will be returned.
func MustCheckPipe(src io.Reader, cmds ...io.ReadWriter) *CmdResult {
	return must1(CheckPipe(src, cmds...))
}

// GetPipe chains commands together and captures the output in a [CmdResult].
func GetPipe(src io.Reader, cmds ...io.ReadWriter) (*CmdResult, error) {
	all := make([]any, len(cmds)+1)
	all[0] = src
	for i, s := range cmds {
		all[i+1] = s
	}
	printCmds(all...)
	r := new(CmdResult)
	dst := new(bytes.Buffer)
	if _, err := Copy(dst, src, cmds...); err != nil {
		if ce := new(Error); errors.As(err, &ce) {
			return nil, ce
		} else if ce := new(CopyError); errors.As(err, &ce) {
			if s, ok := ce.Reader.(io.ReadWriter); ok {
				return nil, NewError(err, s)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	r.Ok = true
	r.Output = strings.Trim(dst.String(), "\n")
	if s, ok := cmds[len(cmds)-1].(io.ReadWriter); ok {
		r.Cmd = s
	}
	return r, nil
}

// MustGetPipe chains commands together and captures the output in a
// [CmdResult].
// It panics if any of the commands fail.
func MustGetPipe(src io.Reader, cmds ...io.ReadWriter) *CmdResult {
	return must1(GetPipe(src, cmds...))
}

func printCmds(a ...any) {
	fmt.Fprintf(Trace, "+ ")
	for i, e := range a {
		if i > 0 {
			fmt.Fprintf(Trace, " | ")
		}
		if str, ok := e.(fmt.Stringer); ok {
			fmt.Fprintf(Trace, str.String())
		} else {
			fmt.Fprintf(Trace, "<stream>")
		}
	}
	fmt.Fprintln(Trace)
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }
