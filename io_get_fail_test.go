package cmdio_test

import (
	"fmt"
	"strings"

	"lesiw.io/cmdio/sys"
)

func Example_getFailure() {
	r, err := sys.Get("ls", "/nonexistent_directory")
	fmt.Println("code:", r.Code)
	fmt.Println("err:", err)
	fmt.Println(strings.Contains(strings.ToLower(r.Log),
		"no such file or directory"))
	// Output:
	// code: 2
	// err: exit status 2
	// true
}
