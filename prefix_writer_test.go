package cmdio

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrefixWriter(t *testing.T) {
	tt := []struct {
		prefix string
		writes []string
		want   string
	}{{
		"> ",
		[]string{""},
		"",
	}, {
		"> ",
		[]string{"test"},
		"> test",
	}, {
		"> ",
		[]string{"Hello\nworld"},
		"> Hello\n> world",
	}, {
		"> ",
		[]string{"Hello\nworld\n"},
		"> Hello\n> world\n",
	}, {
		"> ",
		[]string{"Hello", ", world"},
		"> Hello, world",
	}, {
		"> ",
		[]string{"Hello", ", world\n"},
		"> Hello, world\n",
	}}
	for _, test := range tt {
		t.Run(fmt.Sprintf("%q", test.writes), func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := newPrefixWriter(test.prefix, buf)
			for _, write := range test.writes {
				w.Write([]byte(write))
			}
			if got, want := buf.String(), test.want; got != want {
				t.Errorf("written -want +got\n%s", cmp.Diff(want, got))
			}
		})
	}
}
