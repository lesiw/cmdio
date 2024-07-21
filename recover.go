package cmdio

import (
	"errors"
	"io"
	"os"
)

// Recover recovers from panics produced by Must functions from the cmdio
// package. It is meant to be used in conjunction with Must calls to facilitate
// script-like code.
func Recover(w io.Writer) {
	r := recover()
	if r == nil {
		return
	}
	err, ok := r.(error)
	if !ok {
		panic(r)
	}
	if ce := new(Error); errors.As(err, &ce) {
		ce.Print(w)
		os.Exit(1)
	} else {
		panic(err)
	}
}
