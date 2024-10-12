package cmdio_test

import (
	"io"
	"reflect"
	"testing"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

var runners = map[string]*cmdio.Runner{
	"sys": sys.Runner(),
}

func TestRunners(t *testing.T) {
	cmdio.Trace = io.Discard
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
