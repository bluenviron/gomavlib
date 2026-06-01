package gomavlib

import (
	"io"
)

// Endpoint is an endpoint, which provides Channels.
type Endpoint interface {
	init(node *Node) error
	isEndpoint()
	close()
	oneChannelAtAtime() bool
	isDatagram() bool
	provide() (string, io.ReadWriteCloser, error)
}
