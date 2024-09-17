package sys

import (
	"io"

	"lesiw.io/cmdio"
)

var defaultbox = Box()

func Run(args ...string) error {
	return defaultbox.Run(args...)
}

func MustRun(args ...string) {
	defaultbox.MustRun(args...)
}

func Get(args ...string) (*cmdio.Result, error) {
	return defaultbox.Get(args...)
}

func MustGet(args ...string) *cmdio.Result {
	return defaultbox.MustGet(args...)
}

func Command(args ...string) io.ReadWriter {
	return defaultbox.Command(args...)
}
