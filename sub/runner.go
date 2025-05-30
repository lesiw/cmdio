package sub

import (
	"context"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

type cdr struct {
	rnr *cmdio.Runner
	cmd []string
}

func (c *cdr) Command(
	ctx context.Context, env map[string]string, args ...string,
) cmdio.Command {
	return c.rnr.Commander.Command(ctx, env, append(c.cmd, args...)...)
}

// New instantiates a [cmdio.Runner] that runs subcommands.
func New(cmd ...string) *cmdio.Runner {
	return WithRunner(sys.Runner(), cmd...)
}

// WithRunner instantiates a [cmdio.Runner] that runs subcommands using the
// given runner.
func WithRunner(rnr *cmdio.Runner, cmd ...string) *cmdio.Runner {
	return rnr.WithCommander(&cdr{rnr, cmd})
}
