package sub

import (
	"context"
	"io"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

type cdr struct {
	rnr *cmdio.Runner
	cmd []string
}

func (c *cdr) Command(
	ctx context.Context, env map[string]string, args ...string,
) io.ReadWriter {
	return c.rnr.Commander.Command(ctx, env, append(c.cmd, args...)...)
}

// New instantiates a [cmdio.WithRunner] that runs subcommands.
func New(cmd ...string) *cmdio.Runner {
	return WithRunner(sys.Runner(), cmd...)
}

// WithRunner instantiates a [cmdio.WithRunner] that runs subcommands using the
// given runner.
func WithRunner(rnr *cmdio.Runner, cmd ...string) *cmdio.Runner {
	return cmdio.NewRunner(
		context.Background(),
		make(map[string]string),
		&cdr{rnr, cmd},
	)
}
