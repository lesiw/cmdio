package cmdio

import "io"

type prefixWriter struct {
	p string
	w io.Writer
	n bool
}

func newPrefixWriter(prefix string, w io.Writer) *prefixWriter {
	return &prefixWriter{prefix, w, true}
}

func (w *prefixWriter) Write(p []byte) (n int, err error) {
	var a []byte
	for _, b := range p {
		if w.n {
			a = append(a, []byte(w.p)...)
			w.n = false
		}
		a = append(a, b)
		if b == '\n' {
			w.n = true
		}
	}
	return w.w.Write(a)
}
