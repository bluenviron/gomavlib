package gomavlib

import (
	"io"
)

// EndpointConf is the interface implemented by all endpoint configurations.
type EndpointConf interface {
	init() (Endpoint, error)
}

// Endpoint is an endpoint, which can create Channels.
type Endpoint interface {
	// Conf returns the configuration used to initialize the endpoint
	Conf() EndpointConf
	isEndpoint()
}

// a endpoint must also implement one of the following:
// - endpointChannelSingle
// - endpointChannelAccepter

// endpointChannelSingle is an endpoint that provides a single channel.
// Read() must not return any error unless Close() is called.
type endpointChannelSingle interface {
	Endpoint
	Label() string
	io.ReadWriteCloser
}

// endpointChannelAccepter is an endpoint that provides multiple channels.
type endpointChannelAccepter interface {
	Endpoint
	Close() error
	Accept() (string, io.ReadWriteCloser, error)
}
