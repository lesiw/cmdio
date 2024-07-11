package cmdio

import "io"

// CmdResult describes the results of a completed read from an
// [AttachReadWriter].
type CmdResult struct {
	Ok     bool
	Code   int
	Cmd    io.ReadWriter
	Output string
}
