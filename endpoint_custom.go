package gomavlib

import (
	"io"
)

// EndpointCustom sets up a endpoint that works through a custom interface
// that provides the Read(), Write() and Close() functions.
type EndpointCustom struct {
	// the struct or interface implementing Read(), Write() and Close()
	ReadWriteCloser io.ReadWriteCloser
}

type endpointCustom struct {
	io.ReadWriteCloser
}

func (conf EndpointCustom) init() (endpoint, error) {
	t := &endpointCustom{
		conf.ReadWriteCloser,
	}
	return t, nil
}

func (t *endpointCustom) isEndpoint() {
}
