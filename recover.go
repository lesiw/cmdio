package cmdio

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Recover recovers from panics produced by Must functions from the cmdio
// package. It is meant to be used in conjunction with Must calls to facilitate
// script-like code.
func Recover(w io.Writer) {
	if r := recover(); r != nil {
		err, ok := r.(error)
		if !ok {
			panic(r)
		}
		if ce := new(Error); errors.As(err, &ce) {
			fmt.Fprintf(w, "exec failed: %v: %s\n", ce.Cmd, ce.Error())
			if ce.Log != "" {
				fmt.Fprintf(w, "\nstderr:\n---\n%s\n---\n", ce.Log)
			}
		} else {
			panic(err)
		}
		os.Exit(1)
	}
}
