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
		return nil, err
	}
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
