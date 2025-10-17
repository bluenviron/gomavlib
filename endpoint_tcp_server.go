package gomavlib

import "net"

// EndpointTCPServer sets up a endpoint that works with a TCP server.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (conf EndpointTCPServer) init(node *Node) (Endpoint, error) {
	e := &endpointServer{
		node: node,
		conf: EndpointCustomServer{
			Listen: func() (net.Listener, error) {
				return net.Listen("tcp4", conf.Address)
			},
			Label: "tcp",
		},
	}
	err := e.initialize()
	return e, err
}
