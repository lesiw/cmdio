package sys

import (
	"context"
	"io"

	"lesiw.io/cmdio"
)

type cdr struct{}

func (cdr) Command(
	ctx context.Context, env map[string]string, args ...string,
) io.ReadWriter {
	return newCmd(ctx, env, args...)
}

// Runner instantiates a [cmdio.Runner] that runs commands on the local system.
func Runner() *cmdio.Runner {
	return cmdio.NewRunner(
		context.Background(),
		make(map[string]string),
		new(cdr),
	)
}
