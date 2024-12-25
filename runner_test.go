package cmdio_test

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sub"
	"lesiw.io/cmdio/sys"
	"lesiw.io/prefix"
)

func Example() {
	rnr := sys.Runner().WithEnv(map[string]string{
		"PKGNAME": "cmdio",
	})
	err := rnr.Run("echo", "hello from", rnr.Env("PKGNAME"))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := rnr.Get("true"); err == nil {
		fmt.Println("true always succeeds")
	}
	if _, err := rnr.Get("false"); err != nil {
		fmt.Println("false always fails")
	}
	err = cmdio.Pipe(
		rnr.Command("echo", "pIpEs wOrK tOo"),
		rnr.Command("tr", "A-Z", "a-z"),
	)
	if err != nil {
		log.Fatal(err)
	}
	err = cmdio.Pipe(
		strings.NewReader("Even When Mixed With Other IO"),
		rnr.Command("tr", "A-Z", "a-z"),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// hello from cmdio
	// true always succeeds
	// false always fails
	// pipes work too
	// even when mixed with other io
}

func Example_script() {
	rnr := sys.Runner().WithEnv(map[string]string{
		"PKGNAME": "cmdio",
	})
	rnr.MustRun("echo", "hello from", rnr.Env("PKGNAME"))
	if _, err := rnr.Get("true"); err == nil {
		rnr.MustRun("echo", "true always succeeds")
	}
	if _, err := rnr.Get("false"); err != nil {
		rnr.MustRun("echo", "false always fails")
	}
	cmdio.MustPipe(
		rnr.Command("echo", "pIpEs wOrK tOo"),
		rnr.Command("tr", "A-Z", "a-z"),
	)
	cmdio.MustPipe(
		strings.NewReader("Even When Mixed With Other IO"),
		rnr.Command("tr", "A-Z", "a-z"),
	)
	// Output:
	// hello from cmdio
	// true always succeeds
	// false always fails
	// pipes work too
	// even when mixed with other io
}

func ExampleRunner_Command() {
	rnr := sys.Runner()
	cmd := rnr.Command("echo", "hello world")
	out, err := io.ReadAll(cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	// Output:
	// hello world
}

func ExampleRunner_Run() {
	rnr := sys.Runner()
	err := rnr.Run("echo", "hello world")
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// hello world
}

func ExampleRunner_Get() {
	rnr := sys.Runner()
	r, err := rnr.Get("echo", "hello world")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("out:", r.Out)
	fmt.Println("code:", r.Code)
	// Output:
	// out: hello world
	// code: 0
}

func ExampleRunner_Get_error() {
	rnr := sys.Runner()

	// Use busybox ls to normalize output.
	_, err := rnr.Get("busybox", "ls", "/bad_directory")
	fmt.Println("err:", err)
	// Output:
	// err: exit status 1
	// out: <empty>
	// log:
	// 	ls: /bad_directory: No such file or directory
	// code: 1
}

func ExampleRunner_MustGet() {
	rnr := sys.Runner()
	r := rnr.MustGet("echo", "hello world")
	fmt.Println("out:", r.Out)
	fmt.Println("code:", r.Code)
	// Output:
	// out: hello world
	// code: 0
}

func ExampleRunner_MustGet_panic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	rnr := sys.Runner()

	// Use busybox ls to normalize output.
	rnr.MustGet("busybox", "ls", "/bad_directory")
	// Output:
	// exit status 1
	// out: <empty>
	// log:
	// 	ls: /bad_directory: No such file or directory
	// code: 1
}

func ExampleRunner_WithEnv() {
	rnr := sys.Runner().WithEnv(map[string]string{
		"HOME": "/",
	})
	fmt.Println("rnr(HOME):", rnr.Env("HOME"))
	// Output:
	// rnr(HOME): /
}

func ExampleRunner_WithEnv_multiple() {
	rnr := sys.Runner().WithEnv(map[string]string{
		"HOME": "/",
		"FOO":  "bar",
	})
	fmt.Println("rnr(HOME):", rnr.Env("HOME"))
	rnr2 := rnr.WithEnv(map[string]string{
		"HOME": "/home/example",
	})
	fmt.Println("rnr(HOME):", rnr.Env("HOME"))
	fmt.Println("rnr2(HOME):", rnr2.Env("HOME"))
	fmt.Println("rnr(FOO):", rnr.Env("FOO"))
	fmt.Println("rnr2(FOO):", rnr2.Env("FOO"))
	// Output:
	// rnr(HOME): /
	// rnr(HOME): /
	// rnr2(HOME): /home/example
	// rnr(FOO): bar
	// rnr2(FOO): bar
}

func ExampleRunner_WithEnv_pwd() {
	cmdio.Trace = prefix.NewWriter("+ ", os.Stdout)
	rnr := sys.Runner()

	defer rnr.Run("rm", "-r", "/tmp/cmdio_dir_test")
	err := rnr.Run("mkdir", "/tmp/cmdio_dir_test")
	if err != nil {
		log.Fatal(err)
	}

	err = rnr.WithEnv(map[string]string{
		"PWD": "/tmp/cmdio_dir_test",
	}).Run("pwd")
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// + mkdir /tmp/cmdio_dir_test
	// + PWD=/tmp/cmdio_dir_test pwd
	// /tmp/cmdio_dir_test
	// + rm -r /tmp/cmdio_dir_test
}

func ExampleRunner_WithCommand() {
	rnr := sys.Runner()
	box := sub.WithRunner(rnr, "busybox")
	rnr = rnr.WithCommand("dos2unix", box)
	rnr = rnr.WithCommand("unix2dos", box)
	r, err := cmdio.GetPipe(
		rnr.Command("printf", "hello\r\nworld\r\n"),
		rnr.Command("dos2unix"),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.Out)
	// Output:
	// hello
	// world
}

func ExampleRunner_WithCommand_env() {
	rnr := sys.Runner()
	box := sub.WithRunner(rnr, "busybox")
	rnr = rnr.WithCommand("sh", box)
	err := rnr.WithEnv(map[string]string{"PKGNAME": "cmdio"}).
		Run("sh", "-c", `echo "hello from $PKGNAME"`)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// hello from cmdio
}
