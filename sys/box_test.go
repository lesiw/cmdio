package sys

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"lesiw.io/cmdio"
)

func TestRunSuccess(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		os.Exit(0)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	for i, fn := range []func(...string) error{Run, box.Run} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			err := fn(os.Args[0], "-test.run=TestRunSuccess")
			trcw.WriteString(os.Args[0] + " -test.run=TestRunSuccess\n")
			if err != nil {
				t.Errorf("box.Run() = %q, want <nil>", err)
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Run() trace -want +got\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func TestRunFailure(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		os.Exit(42)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	for i, fn := range []func(...string) error{Run, box.Run} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			err := fn(os.Args[0], "-test.run=TestRunFailure")
			trcw.WriteString(os.Args[0] + " -test.run=TestRunFailure\n")
			if err == nil {
				t.Errorf("box.Run() = <nil>, want exitError")
			} else if ee := new(exec.ExitError); !errors.As(err, &ee) {
				t.Errorf("box.Run() = %T, want exitError", err)
			} else if got, want := ee.ExitCode(), 42; got != want {
				t.Errorf("box.Run().error.ExitCode() = %d, want %d", got, want)
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Run() trace -want +got\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func TestRunBadCommand(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)

	for i, fn := range []func(...string) error{Run, box.Run} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			err := fn("this-command-does-not-exist")
			trcw.WriteString("this-command-does-not-exist\n")

			if err == nil {
				t.Errorf("box.Run() = <nil>, want exec.ErrNotFound")
			} else if !errors.Is(err, exec.ErrNotFound) {
				t.Errorf("box.Run() = %q, want exec.ErrNotFound", err.Error())
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Run() trace -want +got\n%s", cmp.Diff(want, got))
			}
		})
	}
}

func TestGetSuccess(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		fmt.Println("hello world")
		fmt.Fprintln(os.Stderr, "hello stderr")
		os.Exit(0)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	for i, fn := range []func(...string) (*cmdio.Result, error){Get, box.Get} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			r, err := fn(os.Args[0], "-test.run=TestGetSuccess")
			trcw.WriteString(os.Args[0] + " -test.run=TestGetSuccess\n")
			if err != nil {
				t.Errorf("box.Get().error = %q, want <nil>", err)
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Get() trace -want +got\n%s", cmp.Diff(want, got))
			}
			if got, want := r.Out, "hello world"; got != want {
				t.Errorf("box.Get().Result.Out = %q, want %q", got, want)
			}
			if got, want := r.Log, "hello stderr"; got != want {
				t.Errorf("box.Get().Result.Log = %q, want %q", got, want)
			}
			if got, want := r.Code, 0; got != want {
				t.Errorf("box.Get().Result.Code = %d, want %d", got, want)
			}
		})
	}
}

func TestGetFailure(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		fmt.Println("hello world")
		fmt.Fprintln(os.Stderr, "hello stderr")
		os.Exit(42)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	for i, fn := range []func(...string) (*cmdio.Result, error){Get, box.Get} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			r, err := fn(os.Args[0], "-test.run=TestGetFailure")
			trcw.WriteString(os.Args[0] + " -test.run=TestGetFailure\n")
			if err == nil {
				t.Errorf("box.Get().error = <nil>, want error")
			} else if ee := new(exec.ExitError); !errors.As(err, &ee) {
				t.Errorf("box.Get().error = %T, want exitError", err)
			} else if got, want := ee.ExitCode(), 42; got != want {
				t.Errorf("box.Get().error.ExitCode() = %d, want %d", got, want)
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Get() trace -want +got\n%s", cmp.Diff(want, got))
			}
			if r == nil {
				t.Fatalf("box.Get().Result = <nil>, want Result")
			}
			if got, want := r.Out, "hello world"; got != want {
				t.Errorf("box.Get().Result.Out = %q, want %q", got, want)
			}
			if got, want := r.Log, "hello stderr"; got != want {
				t.Errorf("box.Get().Result.Log = %q, want %q", got, want)
			}
			if got, want := r.Code, 42; got != want {
				t.Errorf("box.Get().Result.Code = %d, want %d", got, want)
			}
		})
	}
}

func TestGetBadCommand(t *testing.T) {
	var (
		trc  bytes.Buffer
		trcw bytes.Buffer
		box  = Box()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)

	for i, fn := range []func(...string) (*cmdio.Result, error){Get, box.Get} {
		t.Run(fmt.Sprintf("%T %d", fn, i), func(t *testing.T) {
			r, err := fn("this-command-does-not-exist")
			trcw.WriteString("this-command-does-not-exist\n")
			if err == nil {
				t.Errorf("box.Get().error = <nil>, want error")
			} else if !errors.Is(err, exec.ErrNotFound) {
				t.Errorf("box.Get().error = %q, want exec.ErrNotFound",
					err.Error())
			}
			if got, want := trc.String(), trcw.String(); !cmp.Equal(got, want) {
				t.Errorf("box.Get() trace -want +got\n%s", cmp.Diff(want, got))
			}
			if r == nil {
				t.Fatalf("box.Get().Result = <nil>, want Result")
			}
			if got, want := r.Code, 0; got != want {
				t.Errorf("box.Get().Result.Code = %d, want %d", got, want)
			}
		})
	}
}

func swap[T any](t *testing.T, orig *T, with T) {
	t.Helper()
	o := *orig
	t.Cleanup(func() { *orig = o })
	*orig = with
}
