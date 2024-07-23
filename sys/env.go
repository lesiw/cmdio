package sys

import "lesiw.io/cmdio"

type Enver interface {
	Env(string) string
	Setenv(string, string)
}

func WithEnv(b *cmdio.Box, env map[string]string) *cmdio.Box {
	c, ok := b.Commander.(Enver)
	if !ok {
		panic("bad Commander: not an Enver")
	}
	for k, v := range env {
		c.Setenv(k, v)
	}
	return cmdio.NewBox(c.(cmdio.Commander))
}
