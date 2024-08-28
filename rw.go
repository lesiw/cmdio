package cmdio

import (
	"fmt"
	"io"
)

type nopWriter struct{ io.Reader }
type nopReader struct{ io.Writer }

func (w nopWriter) Write([]byte) (int, error) { return 0, nil }
func (w nopReader) Read([]byte) (int, error)  { return 0, nil }

// readWriter normalizes (io.Reader | io.Writer) to an io.ReadWriter.
// If Go supports union types in the future, this function should be removed.
func readWriter(a any) io.ReadWriter {
	if rw, ok := a.(io.ReadWriter); ok {
		return rw
	} else if r, ok := a.(io.Reader); ok {
		return nopWriter{r}
	} else if w, ok := a.(io.Writer); ok {
		return nopReader{w}
	} else {
		panic(fmt.Sprintf("%T is not an io.Reader or io.Writer", a))
	}
}
