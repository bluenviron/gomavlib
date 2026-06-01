package gomavlib

import "net"

var _ Endpoint = (*EndpointTCPServer)(nil)

// EndpointTCPServer is an endpoint that works with a TCP server.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string

	EndpointCustomServer
}

func (e *EndpointTCPServer) init(node *Node) error {
	e.EndpointCustomServer = EndpointCustomServer{
		Listen: func() (net.Listener, error) {
			return net.Listen("tcp4", e.Address)
		},
		Label: "tcp",
	}
	return e.EndpointCustomServer.init(node)
}
