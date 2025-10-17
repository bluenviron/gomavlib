package gomavlib

import "net"

// EndpointCustomServer sets up a endpoint that works with custom implementations
// by providing a custom Listen func that returns a net.Listener.
// This allows you to use custom protocols that conform to the net.listner.
// A use case could be to add encrypted protocol implementations like DTLS or TCP with TLS.
type EndpointCustomServer struct {
	// function to invoke when server should start listening
	Listen func() (net.Listener, error)

	// the label of the protocol
	Label string
}

func (conf EndpointCustomServer) init(node *Node) (Endpoint, error) {
	e := &endpointServer{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}
