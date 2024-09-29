package cmdio

import (
	"io"

	"golang.org/x/sync/errgroup"
)

// Copy copies the output of each stream into the input of the next stream.
// When output is finished copying from one stream, the receiving stream is
// closed if it is an [io.Closer].
func Copy(
	dst io.Writer, src io.Reader, mid ...io.ReadWriter,
) (written int64, err error) {
	var (
		g errgroup.Group
		r io.Reader
		w io.Writer

		count = make(chan int64)
		total = make(chan int64)
	)
	go func() {
		var written int64
		for n := range count {
			written += n
		}
		total <- written
	}()
	for i := -1; i < len(mid); i++ {
		if i < 0 {
			r = src
		} else {
			r = mid[i]
		}
		if i == len(mid)-1 {
			w = dst
		} else {
			w = mid[i+1]
		}
		i := i
		w := w
		r := r
		g.Go(func() (err error) {
			defer func() {
				if c, ok := w.(io.Closer); ok {
					err1 := c.Close()
					if err == nil {
						err = err1
					}
				}
			}()
			if n, err := io.Copy(w, r); err != nil {
				return copyError{err, i + 1}
			} else {
				count <- n
			}
			return nil
		})
	}
	err = g.Wait()
	close(count)
	return <-total, err
}

type copyError struct {
	err error
	off int
}

func (e copyError) Error() string {
	return e.err.Error()
}

func (e copyError) Unwrap() error {
	return e.err
}

func (e copyError) Offset() int {
	return e.off
}
