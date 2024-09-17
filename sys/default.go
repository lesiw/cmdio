package sys

import (
	"context"
	"io"

	"lesiw.io/cmdio"
)

var box = Box()

func Command(args ...string) io.ReadWriter       { return box.Command(args...) }
func Env(name string) (value string)             { return box.Env(name) }
func Get(args ...string) (*cmdio.Result, error)  { return box.Get(args...) }
func MustGet(args ...string) *cmdio.Result       { return box.MustGet(args...) }
func MustRun(args ...string)                     { box.MustRun(args...) }
func Run(args ...string) error                   { return box.Run(args...) }
func WithContext(ctx context.Context) *cmdio.Box { return box.WithContext(ctx) }
func WithEnv(env map[string]string) *cmdio.Box   { return box.WithEnv(env) }
