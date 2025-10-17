package gomavlib

import (
	"context"
	"net"
)

// EndpointTCPClient sets up a endpoint that works with a TCP client.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (conf EndpointTCPClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: EndpointCustomClient{
			Connect: func(ctx context.Context) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "tcp4", conf.Address)
			},
			Label: "tcp:" + conf.Address,
		},
	}
	err := e.initialize()
	return e, err
}
