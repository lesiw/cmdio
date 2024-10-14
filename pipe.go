package cmdio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

func pipeTrace(src io.Reader, mid []io.ReadWriter) {
	var e any
	for i := -1; i < len(mid); i++ {
		if i < 0 {
			e = src
		} else {
			e = mid[i]
		}
		if i > -1 {
			fmt.Fprintf(Trace, " | ")
		}
		if str, ok := e.(fmt.Stringer); ok {
			fmt.Fprint(Trace, strings.TrimRight(str.String(), "\n"))
		} else {
			fmt.Fprintf(Trace, "<%T>", e)
		}
	}
	fmt.Fprintln(Trace)
}

func pipeErr(src io.Reader, cmd []io.ReadWriter, err error) string {
	var b strings.Builder
	if off, ok := err.(interface{ Offset() int }); ok {
		var e any
		for i := -1; i < len(cmd); i++ {
			if i < 0 {
				e = src
			} else {
				e = cmd[i]
			}
			if i > -1 {
				b.WriteString("\n")
			}
			if str, ok := e.(fmt.Stringer); ok {
				b.WriteString(strings.TrimRight(str.String(), "\n"))
			} else {
				fmt.Fprintf(&b, "<%T>", e)
			}
			if i < len(cmd)-1 {
				b.WriteString(" |")
			}
			if i == off.Offset()-1 {
				b.WriteString(" <- " + err.Error())
			}
		}
	}
	return b.String()
}

// Pipe pipes I/O streams together.
func Pipe(src io.Reader, cmd ...io.ReadWriter) error {
	pipeTrace(src, cmd)
	var e any
	for i := -1; i < len(cmd); i++ {
		if i < 0 {
			e = src
		} else {
			e = cmd[i]
		}
		if l, ok := e.(Logger); ok {
			l.Log(os.Stderr)
		}
	}
	_, err := Copy(nopCloser{os.Stdout}, src, cmd...)
	if err != nil {
		err = fmt.Errorf("%w\n\n%s", err, pipeErr(src, cmd, err))
	}
	return err
}

type nopCloser struct{ io.Writer }

func (nopCloser) Close() error { return nil }

// MustPipe pipes I/O streams together and panics on failure.
func MustPipe(src io.Reader, cmd ...io.ReadWriter) {
	must(Pipe(src, cmd...))
}

// GetPipe pipes I/O streams together and captures the output in a [Result].
func GetPipe(src io.Reader, cmd ...io.ReadWriter) (Result, error) {
	pipeTrace(src, cmd)
	var (
		dst bytes.Buffer
		log syncBuffer
		e   any
		r   Result
	)
	for i := -1; i < len(cmd); i++ {
		if i < 0 {
			e = src
		} else {
			e = cmd[i]
		}
		if l, ok := e.(Logger); ok {
			l.Log(&log)
		}
	}
	_, err := Copy(&dst, src, cmd...)
	r.Out = strings.TrimRight(dst.String(), "\n")
	r.Log = strings.TrimRight(log.String(), "\n")
	r.Cmd = readWriter(e)
	if c, ok := r.Cmd.(Coder); ok {
		r.Code = c.Code()
	}
	if err != nil {
		err = fmt.Errorf("%w\n\n%s\n\nout:%slog:%scode: %d",
			err, pipeErr(src, cmd, err), fmtout(r.Out), fmtout(r.Log), r.Code)
	}
	return r, err
}

// MustGetPipe pipes I/O streams together and captures the output in a
// [Result]. It panics if any of the copy operations fail.
func MustGetPipe(src io.Reader, cmd ...io.ReadWriter) Result {
	return mustv(GetPipe(src, cmd...))
}

type syncBuffer struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func (w *syncBuffer) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *syncBuffer) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}
