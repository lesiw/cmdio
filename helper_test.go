package cmdio

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/google/go-cmp/cmp"
)

func swap[T any](t *testing.T, orig *T, with T) {
	t.Helper()
	o := *orig
	t.Cleanup(func() { *orig = o })
	*orig = with
}

func checkEqual[T any](t *testing.T, name string, got, want T) {
	t.Helper()
	opts := []cmp.Option{
		cmp.Comparer(ptrcmp[bytes.Buffer]),
	}
	if !cmp.Equal(got, want, opts...) {
		t.Errorf("%s -want +got\n%s", name, cmp.Diff(got, want, opts...))
	}
}

func ptrcmp[T any](x, y *T) bool {
	return unsafe.Pointer(x) == unsafe.Pointer(y)
}
