package cmdio_test

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing/iotest"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
	"lesiw.io/prefix"
)

func ExamplePipe_echo() {
	cmdio.Trace = prefix.NewWriter("+ ", os.Stdout)

	rnr := sys.Runner()
	err := cmdio.Pipe(
		rnr.Command("echo", "hello world"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// + echo 'hello world' | tr a-z A-Z
	// HELLO WORLD
}

func ExamplePipe_reader() {
	cmdio.Trace = prefix.NewWriter("+ ", os.Stdout)

	rnr := sys.Runner()
	err := cmdio.Pipe(
		strings.NewReader("hello world"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// + <*strings.Reader> | tr a-z A-Z
	// HELLO WORLD
}

func ExampleMustPipe() {
	rnr := sys.Runner()
	cmdio.MustPipe(
		strings.NewReader("hello world"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	// Output:
	// HELLO WORLD
}

func ExampleMustPipe_panic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	rnr := sys.Runner()
	cmdio.MustPipe(
		iotest.ErrReader(errors.New("some error")),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	// Output:
	// some error
	// <*iotest.errReader> | <- some error
	// tr a-z A-Z
}

func ExampleMustGetPipe() {
	rnr := sys.Runner()
	r := cmdio.MustGetPipe(
		strings.NewReader("hello world"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	fmt.Println("out:", r.Out)
	fmt.Println("code:", r.Code)
	// Output:
	// out: HELLO WORLD
	// code: 0
}

func ExampleMustGetPipe_panic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	rnr := sys.Runner()
	_ = cmdio.MustGetPipe(
		// Use busybox ls to normalize output.
		rnr.Command("busybox", "ls", "/bad_directory"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	// Output:
	// exit status 1
	// busybox ls /bad_directory | <- exit status 1
	// tr a-z A-Z
	//
	// out: <empty>
	// log:
	// 	ls: /bad_directory: No such file or directory
	// code: 0
}
