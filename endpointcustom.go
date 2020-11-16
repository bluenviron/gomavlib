package gomavlib

import (
	"io"
)

// EndpointCustom sets up a endpoint that works with a custom interface
// that provides the Read(), Write() and Close() functions.
type EndpointCustom struct {
	// the struct or interface implementing Read(), Write() and Close()
	ReadWriteCloser io.ReadWriteCloser
}

type endpointCustom struct {
	conf EndpointCustom
	io.ReadWriteCloser
}

func (conf EndpointCustom) init() (Endpoint, error) {
	t := &endpointCustom{
		conf:            conf,
		ReadWriteCloser: conf.ReadWriteCloser,
	}
	return t, nil
}

func (t *endpointCustom) isEndpoint() {}

func (t *endpointCustom) Conf() EndpointConf {
	return t.conf
}

func (t *endpointCustom) Label() string {
	return "custom"
}
