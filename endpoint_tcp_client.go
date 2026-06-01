package gomavlib

import (
	"context"
	"net"
)

var _ Endpoint = (*EndpointTCPClient)(nil)

// EndpointTCPClient is an endpoint that works with a TCP client.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string

	EndpointCustomClient
}

func (e *EndpointTCPClient) init(node *Node) error {
	e.EndpointCustomClient = EndpointCustomClient{
		Connect: func(ctx context.Context) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp4", e.Address)
		},
		Label: "tcp:" + e.Address,
	}
	return e.EndpointCustomClient.init(node)
}
