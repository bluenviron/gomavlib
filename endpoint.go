package gomavlib

import (
	"io"
)

// EndpointConf is the interface implemented by all endpoint configurations.
type EndpointConf interface {
	init() (endpoint, error)
}

// a endpoint must also implement one of the following:
// - endpointChannelSingle
// - endpointChannelAccepter
type endpoint interface{}

// endpoint that provides a single channel.
// Read() must not return any error unless Close() is called.
type endpointChannelSingle interface {
	Label() string
	io.ReadWriteCloser
}

// endpoint that provides multiple channels.
type endpointChannelAccepter interface {
	Close() error
	Accept() (endpointChannelSingle, error)
}
