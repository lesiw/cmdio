package ctr

import (
	"fmt"
	"testing"
)

func TestAlpine(t *testing.T) {
	rnr, err := New("alpine")
	if err != nil {
		t.Fatal(err)
	}
	defer rnr.Close()

	r, err := rnr.Get("which", "apk")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := r.Out, "/sbin/apk"; got != want {
		t.Errorf("[which apk] = %q, want %q", got, want)
	}
}

func TestString(t *testing.T) {
	rnr, err := New("alpine")
	if err != nil {
		t.Fatal(err)
	}
	defer rnr.Close()

	cmd := rnr.Command("echo", "hello world")
	str := fmt.Sprintf("docker container exec %s echo 'hello world'",
		rnr.Commander.(*cdr).ctrid)
	if got, want := fmt.Sprintf("%v", cmd), str; got != want {
		t.Errorf("Sprintf(cmd) = %q, want = %q", got, want)
	}
}

func TestBuild(t *testing.T) {
	rnr, err := New("./testdata/Dockerfile")
	if err != nil {
		t.Fatal(err)
	}
	defer rnr.Close()

	cat, err := rnr.Get("cat", "/tmp/test")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := cat.Out, "hello world"; got != want {
		t.Errorf("[cat /tmp/test] = %q, want %q", got, want)
	}
}
