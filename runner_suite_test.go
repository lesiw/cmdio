package cmdio_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/ctr"
	"lesiw.io/cmdio/sys"
)

var runners = map[string]*cmdio.Runner{
	"sys": sys.Runner(),
	"ctr": mustv(ctr.New("alpine")),
}

func TestRunners(t *testing.T) {
	t.Cleanup(func() {
		for _, rnr := range runners {
			if err := rnr.Close(); err != nil {
				panic(err)
			}
		}
	})
	for name, rnr := range runners {
		suite := reflect.TypeOf(rnrtests{})
		for i := 0; i < suite.NumMethod(); i++ {
			test := suite.Method(i)
			t.Run(name+": "+test.Name, func(t *testing.T) {
				test.Func.Call([]reflect.Value{
					reflect.ValueOf(rnrtests{}),
					reflect.ValueOf(t),
					reflect.ValueOf(rnr),
				})
			})
		}
	}
}

type rnrtests struct{}

func (rnrtests) TestPipe(t *testing.T, rnr *cmdio.Runner) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	rnr = rnr.WithContext(ctx)
	r, err := cmdio.GetPipe(
		rnr.Command("echo", "hello world"),
		rnr.Command("tr", "a-z", "A-Z"),
	)
	if err != nil {
		t.Error(err)
	}
	if got, want := r.Out, "HELLO WORLD"; got != want {
		t.Errorf("[echo hello world] | [tr a-z A-Z] = %q, want %q", got, want)
	}
	cancel()
}

func (rnrtests) TestPipeToCompletedCommand(t *testing.T, rnr *cmdio.Runner) {
	err := cmdio.Pipe(
		rnr.Command("sleep", "0.001"),
		rnr.Command("true"),
	)
	if err != nil {
		t.Errorf("want <nil>, got %q", err.Error())
	}
}

func (rnrtests) TestPipeToFailedCommand(t *testing.T, rnr *cmdio.Runner) {
	err := cmdio.Pipe(
		rnr.Command("sleep", "0.001"),
		rnr.Command("false"),
	)
	if err == nil {
		t.Error("want err, got <nil>")
	}
}

func (rnrtests) TestEnv(t *testing.T, rnr *cmdio.Runner) {
	rnr = rnr.WithEnv(map[string]string{
		"TEST_ENV": "testenv",
	})
	if got, want := rnr.Env("TEST_ENV"), "testenv"; got != want {
		t.Errorf("Env(TEST_ENV) = %q, want %q", got, want)
	}
}

func (rnrtests) TestContext(t *testing.T, rnr *cmdio.Runner) {
	ctx, cancel := context.WithCancel(context.Background())
	rnr = rnr.WithContext(ctx)
	ch := make(chan struct{})
	go func() {
		_, err := rnr.Get("sleep", "5")
		if err == nil {
			t.Error("Get(sleep 5) err = <nil>, want context.Canceled")
		} else if !errors.Is(err, context.Canceled) {
			t.Errorf("Get(sleep 5) err = %q, want context.Canceled", err)
		}
		ch <- struct{}{}
	}()
	cancel()
	<-ch
}

func (rnrtests) TestPwd(t *testing.T, rnr *cmdio.Runner) {
	rnr = rnr.WithEnv(map[string]string{"PWD": "/tmp"})
	r, err := rnr.Get("pwd")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := r.Out, "/tmp"; got != want {
		t.Errorf("[pwd] = %q, want %q", got, want)
	}
}

func mustv[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
