package cmdio

import "io"

// Result describes the results of a command execution.
type Result struct {
	Cmd  io.ReadWriter
	Out  string
	Log  string
	Code int
}
