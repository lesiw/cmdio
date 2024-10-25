package sys

import (
	"context"

	"lesiw.io/cmdio"
)

type cdr struct{}

func (cdr) Command(
	ctx context.Context, env map[string]string, args ...string,
) cmdio.Command {
	return newCmd(ctx, env, args...)
}

// Runner instantiates a [cmdio.Runner] that runs commands on the local system.
func Runner() *cmdio.Runner {
	return new(cmdio.Runner).
		WithCommander(new(cdr)).
		WithContext(context.Background())
}
