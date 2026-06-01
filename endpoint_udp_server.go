package gomavlib

import (
	"net"

	"github.com/pion/transport/v2/udp"
)

var _ Endpoint = (*EndpointUDPServer)(nil)

// EndpointUDPServer is an endpoint that works with an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GCS
// if they are connected to the same network.
type EndpointUDPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string

	EndpointCustomServer
}

func (e *EndpointUDPServer) init(node *Node) error {
	e.EndpointCustomServer = EndpointCustomServer{
		Listen: func() (net.Listener, error) {
			addr, err := net.ResolveUDPAddr("udp4", e.Address)
			if err != nil {
				return nil, err
			}

			return udp.Listen("udp4", addr)
		},
		Label:      "udp:" + e.Address,
		IsDatagram: true,
	}
	return e.EndpointCustomServer.init(node)
}
