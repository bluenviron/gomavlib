package gomavlib

import (
	"io"
)

// EndpointCustom sets up a endpoint that works with a custom interface
// that provides the Read(), Write() and Close() functions.
type EndpointCustom struct {
	// struct or interface implementing Read(), Write() and Close()
	ReadWriteCloser io.ReadWriteCloser
}

func (conf EndpointCustom) init(node *Node) (Endpoint, error) {
	e := &endpointCustom{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

type endpointCustom struct {
	node *Node
	conf EndpointCustom

	io.ReadWriteCloser
}

func (e *endpointCustom) initialize() error {
	e.ReadWriteCloser = e.conf.ReadWriteCloser
	return nil
}

func (e *endpointCustom) isEndpoint() {}

func (e *endpointCustom) Conf() EndpointConf {
	return e.conf
}

func (e *endpointCustom) label() string {
	return "custom"
}
