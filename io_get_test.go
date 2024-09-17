package cmdio_test

import (
	"fmt"

	"lesiw.io/cmdio/sys"
)

func Example_get() {
	r, err := sys.Get("echo", "hello world")
	if err != nil {
		panic(err)
	}
	fmt.Println("code:", r.Code)
	fmt.Println("out:", r.Out)
	// Output:
	// code: 0
	// out: hello world
}
