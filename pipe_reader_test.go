package cmdio_test

import (
	"os"
	"strings"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

func Example_pipeReader() {
	defer cmdio.Recover(os.Stderr)

	cmdio.MustPipe(
		strings.NewReader("hello world"),
		sys.Command("tr", "a-z", "A-Z"),
	)
	// Output:
	// HELLO WORLD
}
