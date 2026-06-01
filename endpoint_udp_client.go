package gomavlib

import (
	"context"
	"net"
)

var _ Endpoint = (*EndpointUDPClient)(nil)

// EndpointUDPClient is an endpoint that works with a UDP client.
type EndpointUDPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string

	EndpointCustomClient
}

func (e *EndpointUDPClient) init(node *Node) error {
	e.EndpointCustomClient = EndpointCustomClient{
		Connect: func(ctx context.Context) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "udp4", e.Address)
		},
		Label:      "udp:" + e.Address,
		IsDatagram: true,
	}
	return e.EndpointCustomClient.init(node)
}
