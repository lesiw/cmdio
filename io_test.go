package cmdio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"

	"github.com/google/go-cmp/cmp"
)

func TestRun(t *testing.T) {
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Trace.(*prefixWriter).w, errbuf)
	cmd := bytes.NewBufferString("hello world")

	err := Run(cmd)

	if err != nil {
		t.Errorf("Run(%q) = %q, want <nil>", cmd, err)
	}
	if got, want := outbuf.String(), "hello world"; got != want {
		t.Errorf("Run(%q) stdout = %q, want %q", cmd, got, want)
	}
	if got, want := errbuf.String(), "+ hello world\n"; got != want {
		t.Errorf("Run(%q) stderr = %q, want %q", cmd, got, want)
	}
	// Ensure nothing was written to cmd after being read.
	if got, want := cmd.String(), ""; got != want {
		t.Errorf("cmd = %q, want %q", got, want)
	}
}

func TestRunError(t *testing.T) {
	swap(t, &Trace, io.Discard)
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Stderr, errbuf)
	cmd := iotest.ErrReader(errors.New("some error"))

	err := Run(cmd)

	if got, want := err.Error(), "some error"; got != want {
		t.Errorf("Run() = %q, want %q", got, want)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Run() stdout = %q, want %q", got, want)
	}
	if got, want := errbuf.String(), ""; got != want {
		t.Errorf("Run() stderr = %q, want %q", got, want)
	}
}

type testAttacher struct {
	attaches    int
	reads       [][]byte
	attachError error
	readError   error
}

func (a *testAttacher) Attach() error {
	a.attaches++
	return a.attachError
}

func (a *testAttacher) Read(p []byte) (int, error) {
	a.reads = append(a.reads, p)
	return 0, a.readError
}

func TestRunAttacher(t *testing.T) {
	swap(t, &Trace, io.Discard)
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Stderr, errbuf)
	a := new(testAttacher)

	err := Run(a)

	if err != nil {
		t.Errorf("Run() = %q, want <nil>", err)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Run() stdout = %q, want %q", got, want)
	}
	if got, want := errbuf.String(), ""; got != want {
		t.Errorf("Run() stderr = %q, want %q", got, want)
	}
	if got, want := a.attaches, 1; got != want {
		t.Errorf("testAttacher.attaches = %d, want %d", got, want)
	}
	if got, want := len(a.reads), 1; got != want {
		t.Errorf("len(testAttacher.reads) = %d, want %d", got, want)
	}
	if got, want := a.reads, [][]byte{nil}; !cmp.Equal(got, want) {
		t.Errorf("testAttacher.reads -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestRunAttacherAttachError(t *testing.T) {
	swap(t, &Trace, io.Discard)
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Stderr, errbuf)
	a := new(testAttacher)
	a.attachError = errors.New("attach error")

	err := Run(a)

	if got, want := err.Error(), "attach error"; got != want {
		t.Errorf("Run() = %q, want %q", got, want)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Run() stdout = %q, want %q", got, want)
	}
	if got, want := errbuf.String(), ""; got != want {
		t.Errorf("Run() stderr = %q, want %q", got, want)
	}
	if got, want := a.attaches, 1; got != want {
		t.Errorf("testAttacher.attaches = %d, want %d", got, want)
	}
	if got, want := len(a.reads), 0; got != want {
		t.Errorf("len(testAttacher.reads) = %d, want %d", got, want)
	}
}

func TestRunAttacherReadError(t *testing.T) {
	swap(t, &Trace, io.Discard)
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Stderr, errbuf)
	a := new(testAttacher)
	a.readError = errors.New("read error")

	err := Run(a)

	if got, want := err.Error(), "read error"; got != want {
		t.Errorf("Run() = %q, want %q", got, want)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Run() stdout = %q, want %q", got, want)
	}
	if got, want := errbuf.String(), ""; got != want {
		t.Errorf("Run() stderr = %q, want %q", got, want)
	}
	if got, want := a.attaches, 1; got != want {
		t.Errorf("testAttacher.attaches = %d, want %d", got, want)
	}
	if got, want := len(a.reads), 1; got != want {
		t.Errorf("len(testAttacher.reads) = %d, want %d", got, want)
	}
	if got, want := a.reads, [][]byte{nil}; !cmp.Equal(got, want) {
		t.Errorf("testAttacher.reads -want +got\n%s", cmp.Diff(want, got))
	}
}

func TestGet(t *testing.T) {
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Trace.(*prefixWriter).w, errbuf)
	cmd := bytes.NewBufferString("hello world")

	r, err := Get(cmd)

	checkEqual(t, fmt.Sprintf("Get(%q).CmdResult", cmd), r, &Result{
		Cmd: cmd, Out: "hello world",
	})
	if err != nil {
		t.Errorf("Get(%q).error = %q, want <nil>", cmd, err)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Get(%q) stdout = %q, want %q", cmd, got, want)
	}
	if got, want := errbuf.String(), "+ hello world\n"; got != want {
		t.Errorf("Run(%q) stderr = %q, want %q", cmd, got, want)
	}
	// Ensure nothing was written to cmd after being read.
	if got, want := cmd.String(), ""; got != want {
		t.Errorf("cmd = %q, want %q", got, want)
	}
}

func TestGetError(t *testing.T) {
	swap(t, &Trace, io.Discard)
	outbuf, errbuf := new(bytes.Buffer), new(bytes.Buffer)
	swap[io.Writer](t, &Stdout, outbuf)
	swap[io.Writer](t, &Stderr, errbuf)
	cmd := iotest.ErrReader(errors.New("some error"))

	r, err := Get(cmd)

	if r == nil {
		t.Error("Get().CmdResult = <nil>, want Result")
	}
	if got, want := err.Error(), "some error"; got != want {
		t.Errorf("Get().error = %q, want %q", got, want)
	}
	if got, want := outbuf.String(), ""; got != want {
		t.Errorf("Get() stdout = %q, want %q", got, want)
	}
	if got, want := errbuf.String(), ""; got != want {
		t.Errorf("Get() stderr = %q, want %q", got, want)
	}
}
