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
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		os.Exit(0)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	err := run.Run(os.Args[0], "-test.run=TestRunSuccess")
	log.WriteString(os.Args[0] + " -test.run=TestRunSuccess\n")
	if err != nil {
		t.Errorf("rnr.Run() = %q, want <nil>", err)
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Run() trace -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestRunFailure(t *testing.T) {
	var (
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		os.Exit(42)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	err := run.Run(os.Args[0], "-test.run=TestRunFailure")
	log.WriteString(os.Args[0] + " -test.run=TestRunFailure\n")
	if err == nil {
		t.Errorf("rnr.Run() = <nil>, want exitError")
	} else if ee := new(exec.ExitError); !errors.As(err, &ee) {
		t.Errorf("rnr.Run() = %T, want exitError", err)
	} else if got, want := ee.ExitCode(), 42; got != want {
		t.Errorf("rnr.Run().error.ExitCode() = %d, want %d", got, want)
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Run() trace -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestRunBadCommand(t *testing.T) {
	var (
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)

	err := run.Run("this-command-does-not-exist")
	log.WriteString("this-command-does-not-exist\n")

	if err == nil {
		t.Errorf("rnr.Run() = <nil>, want exec.ErrNotFound")
	} else if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("rnr.Run() = %q, want exec.ErrNotFound", err.Error())
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Run() trace -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestGetSuccess(t *testing.T) {
	var (
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		fmt.Println("hello world")
		fmt.Fprintln(os.Stderr, "hello stderr")
		os.Exit(0)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	r, err := run.Get(os.Args[0], "-test.run=TestGetSuccess")
	log.WriteString(os.Args[0] + " -test.run=TestGetSuccess\n")
	if err != nil {
		t.Errorf("rnr.Get().error = %q, want <nil>", err)
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Get() trace -want +got\n%s", cmp.Diff(want, got))
	}
	if got, want := r.Out, "hello world"; got != want {
		t.Errorf("rnr.Get().Result.Out = %q, want %q", got, want)
	}
	if got, want := r.Log, "hello stderr"; got != want {
		t.Errorf("rnr.Get().Result.Log = %q, want %q", got, want)
	}
	if got, want := r.Code, 0; got != want {
		t.Errorf("rnr.Get().Result.Code = %d, want %d", got, want)
	}
}

func TestGetFailure(t *testing.T) {
	var (
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)
	if os.Getenv("CMD_TEST_PROC") == "1" {
		fmt.Println("hello world")
		fmt.Fprintln(os.Stderr, "hello stderr")
		os.Exit(42)
	}
	t.Setenv("CMD_TEST_PROC", "1")

	r, err := run.Get(os.Args[0], "-test.run=TestGetFailure")
	log.WriteString(os.Args[0] + " -test.run=TestGetFailure\n")
	if err == nil {
		t.Errorf("rnr.Get().error = <nil>, want error")
	} else if ee := new(exec.ExitError); !errors.As(err, &ee) {
		t.Errorf("rnr.Get().error = %T, want exitError", err)
	} else if got, want := ee.ExitCode(), 42; got != want {
		t.Errorf("rnr.Get().error.ExitCode() = %d, want %d", got, want)
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Get() trace -want +got\n%s", cmp.Diff(want, got))
	}
	if got, want := r.Out, "hello world"; got != want {
		t.Errorf("rnr.Get().Result.Out = %q, want %q", got, want)
	}
	if got, want := r.Log, "hello stderr"; got != want {
		t.Errorf("rnr.Get().Result.Log = %q, want %q", got, want)
	}
	if got, want := r.Code, 42; got != want {
		t.Errorf("rnr.Get().Result.Code = %d, want %d", got, want)
	}
}

func TestGetBadCommand(t *testing.T) {
	var (
		trc bytes.Buffer
		log bytes.Buffer
		run = Runner()
	)
	swap[io.Writer](t, &cmdio.Trace, &trc)

	r, err := run.Get("this-command-does-not-exist")
	log.WriteString("this-command-does-not-exist\n")
	if err == nil {
		t.Errorf("rnr.Get().error = <nil>, want error")
	} else if !errors.Is(err, exec.ErrNotFound) {
		t.Errorf("rnr.Get().error = %q, want exec.ErrNotFound",
			err.Error())
	}
	if got, want := trc.String(), log.String(); !cmp.Equal(got, want) {
		t.Errorf("rnr.Get() trace -want +got\n%s", cmp.Diff(want, got))
	}
	if got, want := r.Code, 0; got != want {
		t.Errorf("rnr.Get().Result.Code = %d, want %d", got, want)
	}
}

func swap[T any](t *testing.T, orig *T, with T) {
	t.Helper()
	o := *orig
	t.Cleanup(func() { *orig = o })
	*orig = with
}
