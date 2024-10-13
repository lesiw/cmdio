package ctr

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sub"
	"lesiw.io/cmdio/sys"
)

var clis = [...][]string{
	{"docker"},
	{"podman"},
	{"nerdctl"},
	{"lima", "nerdctl"},
}

type cdr struct {
	rnr   *cmdio.Runner
	ctrid string
}

func (c *cdr) Command(
	ctx context.Context, env map[string]string, args ...string,
) io.ReadWriter {
	return newCmd(c, ctx, env, args...)
}

func (c *cdr) Close() error {
	return c.rnr.Run("container", "rm", "-f", c.ctrid)
}

// New instantiates a [cmdio.Runner] that runs commands in a container.
func New(name string) (*cmdio.Runner, error) {
	return WithRunner(sys.Runner(), name)
}

// WithRunner instantiates a [cmdio.Runner] that runs commands in a container
// using the given runner.
func WithRunner(rnr *cmdio.Runner, name string) (*cmdio.Runner, error) {
	var ctrcli []string
	for _, cli := range clis {
		if _, err := rnr.Get("which", cli[0]); err == nil {
			ctrcli = cli
			break
		}
	}
	if len(ctrcli) == 0 {
		return nil, fmt.Errorf("failed to find container CLI")
	}
	rnr = sub.WithRunner(rnr, ctrcli...)

	if len(name) > 0 && (name[0] == '/' || name[0] == '.') {
		var err error
		if name, err = buildContainer(rnr, name); err != nil {
			return nil, fmt.Errorf("failed to build container: %w", err)
		}
	}
	r, err := rnr.Get("container", "run", "--rm", "-d", "-i", name, "cat")
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	return cmdio.NewRunner(
		context.Background(),
		make(map[string]string),
		&cdr{rnr: rnr, ctrid: r.Out},
	), nil
}

func buildContainer(
	rnr *cmdio.Runner, rpath string,
) (image string, err error) {
	var path string
	if path, err = filepath.Abs(rpath); err != nil {
		err = fmt.Errorf("bad Containerfile path '%s': %w", rpath, err)
		return
	}
	imagehash := sha1.New()
	imagehash.Write([]byte(path))
	image = fmt.Sprintf("%x", imagehash.Sum(nil))
	insp, insperr := rnr.Get(
		"image", "inspect",
		"--format", "{{.Created}}",
		image,
	)
	mtime, err := getMtime(path)
	if err != nil {
		err = fmt.Errorf("bad Containerfile '%s': %w", path, err)
		return
	}
	if insperr == nil {
		var ctime time.Time
		ctime, err = time.Parse(time.RFC3339, insp.Out)
		if err != nil {
			err = fmt.Errorf(
				"failed to parse container created timestamp '%s': %s",
				insp.Out, err)
			return
		}
		if ctime.Unix() > mtime {
			return // Container is newer than Containerfile.
		}
	}
	err = rnr.Run(
		"image", "build",
		"--file", path,
		"--no-cache",
		"--tag", image,
		filepath.Dir(path),
	)
	if err != nil {
		err = fmt.Errorf("failed to build '%s': %w", path, err)
	}
	return
}

func getMtime(path string) (mtime int64, err error) {
	var info fs.FileInfo
	info, err = os.Lstat(path)
	if err != nil {
		return
	}
	mtime = info.ModTime().Unix()
	return
}
