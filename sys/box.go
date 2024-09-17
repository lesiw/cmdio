package sys

import (
	"context"
	"io"

	"lesiw.io/cmdio"
)

type sysbox struct{}

func (b *sysbox) Command(
	ctx context.Context, env map[string]string, args ...string,
) io.ReadWriter {
	return newCmd(ctx, env, args...)
}

func Box() *cmdio.Box {
	return cmdio.NewBox(
		context.Background(),
		make(map[string]string),
		new(sysbox),
	)
}
