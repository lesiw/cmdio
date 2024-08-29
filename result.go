package cmdio

import "io"

// CmdResult describes the results of a completed read from an [io.ReadWriter].
type CmdResult struct {
	Cmd io.ReadWriter
	Out string
}
