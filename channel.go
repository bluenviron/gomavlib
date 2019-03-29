package gomavlib

import (
	"io"
)

// Channel is a channel provided by a endpoint.
type Channel struct {
	// the endpoint which the channel belongs to
	Endpoint Endpoint

	label     string
	rwc       io.ReadWriteCloser
	writeChan chan interface{}
}

// String implements fmt.Stringer and returns the channel label.
func (e *Channel) String() string {
	return e.label
}
