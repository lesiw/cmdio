package cmdio_test

import (
	"lesiw.io/cmdio/sys"
)

func Example_run() {
	if err := sys.Run("echo", "hello world"); err != nil {
		panic(err)
	}
	// Output:
	// hello world
}
