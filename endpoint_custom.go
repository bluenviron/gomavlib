package gomavlib

import (
	"io"
)

type removeCloser struct {
	wrapped io.ReadWriteCloser
}

func (r *removeCloser) Read(p []byte) (int, error) {
	return r.wrapped.Read(p)
}

func (r *removeCloser) Write(p []byte) (int, error) {
	return r.wrapped.Write(p)
}

func (r *removeCloser) Close() error {
	return nil
}

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

	rwc io.ReadWriteCloser
}

func (e *endpointCustom) close() {
	e.rwc.Close()
}

func (e *endpointCustom) initialize() error {
	e.rwc = e.conf.ReadWriteCloser
	return nil
}

func (e *endpointCustom) isEndpoint() {}

func (e *endpointCustom) Conf() EndpointConf {
	return e.conf
}

func (e *endpointCustom) oneChannelAtAtime() bool {
	return true
}

func (e *endpointCustom) provide() (string, io.ReadWriteCloser, error) {
	return "custom", &removeCloser{e.rwc}, nil
}
