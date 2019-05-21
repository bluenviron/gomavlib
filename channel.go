package gomavlib

import (
	"io"
)

// Channel is a communication channel created by an endpoint. For instance, a
// TCP client endpoint creates a single channel, while a TCP server endpoint
// creates a channel for each incoming connection.
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
