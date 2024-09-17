package cmdio_test

import (
	"os"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

func Example_env() {
	defer cmdio.Recover(os.Stderr)

	box := sys.WithEnv(map[string]string{"HOME": "/"})
	box.MustRun("echo", box.Env("HOME"))
	// Output:
	// /
}
