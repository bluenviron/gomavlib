package gomavlib

import (
	"context"
	"net"
)

// EndpointCustomClient sets up a endpoint that works with a custom implementation
// by providing a Connect func that returns a net.Conn.
type EndpointCustomClient struct {
	// custom connect function that opens the connection
	Connect func(ctx context.Context) (net.Conn, error)

	// the label of the protocol
	Label string
}

func (conf EndpointCustomClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}
