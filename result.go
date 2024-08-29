package cmdio

import "io"

// Result describes a completed read from an [io.ReadWriter].
type Result struct {
	Cmd io.ReadWriter
	Out string
}
