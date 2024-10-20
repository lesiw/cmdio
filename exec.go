package cmdio

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/sync/errgroup"
)

func run(cmd io.Reader) error {
	a, ok := cmd.(Attacher)
	if !ok {
		// If this command does not implement Attacher, stream it to stdout
		// (and stderr, if applicable).
		if l, ok := cmd.(Logger); ok {
			l.Log(stderr)
		}
		fmt.Fprintln(Trace, strings.TrimRight(fmt.Sprintf("%v", cmd), "\n"))
		_, err := io.Copy(stdout, cmd)
		return err
	}
	if err := a.Attach(); err != nil {
		return err
	}
	fmt.Fprintln(Trace, strings.TrimRight(fmt.Sprintf("%v", cmd), "\n"))
	_, err := cmd.Read(nil)
	if err == io.EOF {
		err = nil
	}
	return err
}

func get(cmd io.Reader) (Result, error) {
	fmt.Fprintln(Trace, strings.TrimRight(fmt.Sprintf("%v", cmd), "\n"))

	var r Result
	var wg errgroup.Group
	var log bytes.Buffer
	out := make(chan string)

	if l, ok := cmd.(Logger); ok {
		l.Log(&log)
	}
	wg.Go(func() error {
		buf, err := io.ReadAll(cmd)
		out <- strings.TrimRight(string(buf), "\n")
		return err
	})

	r.Cmd = readWriter(cmd)
	r.Out = <-out
	err := wg.Wait()
	r.Log = strings.TrimRight(log.String(), "\n")
	if c, ok := cmd.(Coder); ok {
		r.Code = c.Code()
	}

	return r, err
}
