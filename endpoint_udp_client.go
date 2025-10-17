package gomavlib

import (
	"context"
	"net"
)

// EndpointUDPClient sets up a endpoint that works with a UDP client.
type EndpointUDPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (conf EndpointUDPClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: EndpointCustomClient{
			Connect: func(ctx context.Context) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "udp4", conf.Address)
			},
			Label: "udp:" + conf.Address,
		},
	}
	err := e.initialize()
	return e, err
}
