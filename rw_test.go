package cmdio

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestReadWriterNop(t *testing.T) {
	buf := new(bytes.Buffer)
	rw := readWriter(buf)
	if buf2, ok := rw.(*bytes.Buffer); !ok || !ptrcmp(buf, buf2) {
		t.Errorf("readWriter(rw) != rw")
	}
}

func TestReadWriterReader(t *testing.T) {
	r := strings.NewReader("hello world")
	rw := readWriter(r)
	buf, err := io.ReadAll(rw)
	if err != nil {
		t.Fatalf("io.ReadAll(%+v).error = %q, want <nil>", rw, err)
	}
	if got, want := string(buf), "hello world"; got != want {
		t.Errorf("io.ReadAll(%+v).[]bytes = %q, want %q", rw, got, want)
	}
}

func TestReadWriterWriter(t *testing.T) {
	buf := new(bytes.Buffer)
	var w io.Writer = buf
	rw := readWriter(w)
	_, err := rw.Write([]byte("hello world"))
	if err != nil {
		t.Fatalf("rw.Write().error = %q, want <nil>", err)
	}
	n, err := rw.Read(nil)
	if got, want := n, 0; got != want {
		t.Errorf("rw.Read().n = %d, want %d", got, want)
	}
	if err != nil {
		t.Errorf("rw.Read().error = %q, want <nil>", err)
	}
	if got, want := buf.String(), "hello world"; got != want {
		t.Errorf("buf.String() = %q, want %q", got, want)
	}
}

func TestReadWriterInvalid(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("readWriter() did not panic")
		}
		s, ok := r.(string)
		if !ok {
			t.Fatalf("readWriter() panic was not string")
		}
		want := "struct {} is not an io.Reader or io.Writer"
		if s != want {
			t.Errorf("readWriter() panic = %q, want %q", s, want)
		}
	}()
	readWriter(struct{}{})
}
