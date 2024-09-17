package cmdio_test

import (
	"os"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

func Example_pipeEcho() {
	defer cmdio.Recover(os.Stderr)

	cmdio.MustPipe(
		sys.Command("echo", "hello world"),
		sys.Command("tr", "a-z", "A-Z"),
	)
	// Output:
	// HELLO WORLD
}
