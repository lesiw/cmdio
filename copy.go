package cmdio

import (
	"io"

	"golang.org/x/sync/errgroup"
)

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
				return err
			} else {
				count <- n
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return 0, err
	}
	close(count)
	return <-total, nil
}
