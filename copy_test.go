package cmdio_test

import (
	"io"
	"log"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

// nolint: errcheck
func ExampleCopy() {
	rnr := sys.Runner()

	defer rnr.Run("rm", "-f", "/tmp/cmdio_test.txt")
	_, err := cmdio.Copy(
		io.Discard,
		rnr.Command("echo", "hello world"),
		rnr.Command("tee", "/tmp/cmdio_test.txt"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.Run("cat", "/tmp/cmdio_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// hello world
}
