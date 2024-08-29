package cmdio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Pipe chains commands together.
func Pipe(src io.Reader, cmds ...io.ReadWriter) error {
	all := make([]any, len(cmds)+1)
	all[0] = src
	for i, s := range cmds {
		all[i+1] = s
	}
	printCmds(all...)
	_, err := Copy(nopCloser{os.Stdout}, src, cmds...)
	return err
}

// MustPipe chains commands together and panics on failure.
func MustPipe(src io.Reader, cmds ...io.ReadWriter) {
	must(Pipe(src, cmds...))
}

// GetPipe chains commands together and captures the output in a [Result].
func GetPipe(src io.Reader, cmds ...io.ReadWriter) (*Result, error) {
	all := make([]any, len(cmds)+1)
	all[0] = src
	for i, s := range cmds {
		all[i+1] = s
	}
	printCmds(all...)
	r := new(Result)
	dst := new(bytes.Buffer)
	if _, err := Copy(dst, src, cmds...); err != nil {
		return nil, err
	}
	r.Out = strings.Trim(dst.String(), "\n")
	r.Cmd = cmds[len(cmds)-1]
	if l, ok := r.Cmd.(Logger); ok {
		logbuf, err := io.ReadAll(l.Log())
		if err != nil {
			return nil, fmt.Errorf("read from cmd log failed: %w", err)
		}
		r.Log = strings.Trim(string(logbuf), "\n")
	}
	if c, ok := r.Cmd.(Coder); ok {
		r.Code = c.Code()
	}
	return r, nil
}

// MustGetPipe chains commands together and captures the output in a [Result].
// It panics if any of the commands fail.
func MustGetPipe(src io.Reader, cmds ...io.ReadWriter) *Result {
	return must1(GetPipe(src, cmds...))
}

func printCmds(a ...any) {
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
