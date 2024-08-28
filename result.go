package cmdio

import "io"

// CmdResult describes the results of a completed read from an [io.ReadWriter].
// TODO: Incorporate stderr.
type CmdResult struct {
	Cmd    io.ReadWriter
	Output string
}
