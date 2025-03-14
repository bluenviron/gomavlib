package gomavlib

import (
	"io"
)

// Endpoint is an endpoint, which provides Channels.
type Endpoint interface {
	// Conf returns the configuration used to initialize the endpoint
	Conf() EndpointConf

	isEndpoint()
	close()
	oneChannelAtAtime() bool
	provide() (string, io.ReadWriteCloser, error)
}
