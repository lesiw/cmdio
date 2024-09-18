package cmdio_test

import (
	"os"

	"lesiw.io/cmdio"
	"lesiw.io/cmdio/sys"
)

func Example_chdir() {
	defer cmdio.Recover(os.Stderr)

	defer sys.Run("rm", "-r", "/tmp/cmdio_dir_test")
	sys.MustRun("mkdir", "/tmp/cmdio_dir_test")
	sys.WithEnv(map[string]string{
		"PWD": "/tmp/cmdio_dir_test",
	}).MustRun("pwd")
	// Output: /tmp/cmdio_dir_test
}
