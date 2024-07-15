package cmdio

import "io"

// CmdResult describes the results of a completed read from an
// [AttachReadWriter].
type CmdResult struct {
	Cmd    io.ReadWriter
	Output string
}
