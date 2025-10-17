package gomavlib

import (
	"net"

	"github.com/pion/transport/v2/udp"
)

// EndpointUDPServer sets up a endpoint that works with an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GCS
// if they are connected to the same network.
type EndpointUDPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (conf EndpointUDPServer) init(node *Node) (Endpoint, error) {
	e := &endpointServer{
		node: node,
		conf: EndpointCustomServer{
			Listen: func() (net.Listener, error) {
				addr, err := net.ResolveUDPAddr("udp4", conf.Address)
				if err != nil {
					return nil, err
				}

				return udp.Listen("udp4", addr)
			},
			Label: "udp",
		},
	}
	err := e.initialize()
	return e, err
}
