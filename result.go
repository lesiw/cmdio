package cmdio

import "io"

// CmdResult describes the results of a completed read from an [io.Reader].
type CmdResult struct {
	Cmd    io.Reader
	Output string
}
