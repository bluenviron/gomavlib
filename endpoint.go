package gomavlib

import (
	"io"
)

// EndpointConf is the interface implemented by all endpoint configurations.
type EndpointConf interface {
	init() (Endpoint, error)
}

// Endpoint represents an endpoint, which contains zero or more channels
type Endpoint interface {
	// Conf returns the configuration used to initialize the endpoint
	Conf() interface{}

	isEndpoint()
}

// a endpoint must also implement one of the following:
// - endpointChannelSingle
// - endpointChannelAccepter

// endpoint that provides a single channel.
// Read() must not return any error unless Close() is called.
type endpointChannelSingle interface {
	Endpoint
	Label() string
	io.ReadWriteCloser
}

// endpoint that provides multiple channels.
type endpointChannelAccepter interface {
	Endpoint
	Close() error
	Accept() (string, io.ReadWriteCloser, error)
}
