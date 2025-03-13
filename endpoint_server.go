package gomavlib

import (
	"fmt"
	"io"
	"net"

	"github.com/pion/transport/v2/udp"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

type endpointServerConf interface {
	isUDP() bool
	getAddress() string
	init(*Node) (Endpoint, error)
}

// EndpointTCPServer sets up a endpoint that works with a TCP server.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointTCPServer) isUDP() bool {
	return false
}

func (conf EndpointTCPServer) getAddress() string {
	return conf.Address
}

func (conf EndpointTCPServer) init(node *Node) (Endpoint, error) {
	e := &endpointServer{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

// EndpointUDPServer sets up a endpoint that works with an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GCS
// if they are connected to the same network.
type EndpointUDPServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointUDPServer) isUDP() bool {
	return true
}

func (conf EndpointUDPServer) getAddress() string {
	return conf.Address
}

func (conf EndpointUDPServer) init(node *Node) (Endpoint, error) {
	e := &endpointServer{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

type endpointServer struct {
	node *Node
	conf endpointServerConf

	listener net.Listener

	// in
	terminate chan struct{}
}

func (e *endpointServer) initialize() error {
	_, _, err := net.SplitHostPort(e.conf.getAddress())
	if err != nil {
		return fmt.Errorf("invalid address")
	}

	if e.conf.isUDP() {
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp4", e.conf.getAddress())
		if err != nil {
			return err
		}

		e.listener, err = udp.Listen("udp4", addr)
		if err != nil {
			return err
		}
	} else {
		e.listener, err = net.Listen("tcp4", e.conf.getAddress())
		if err != nil {
			return err
		}
	}

	e.terminate = make(chan struct{})

	return nil
}

func (e *endpointServer) isEndpoint() {}

func (e *endpointServer) Conf() EndpointConf {
	return e.conf
}

func (e *endpointServer) close() {
	close(e.terminate)
	e.listener.Close()
}

func (e *endpointServer) oneChannelAtAtime() bool {
	return false
}

func (e *endpointServer) provide() (string, io.ReadWriteCloser, error) {
	nconn, err := e.listener.Accept()
	// wait termination, do not report errors
	if err != nil {
		<-e.terminate
		return "", nil, errTerminated
	}

	label := fmt.Sprintf("%s:%s", func() string {
		if e.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), nconn.RemoteAddr())

	conn := timednetconn.New(
		e.node.IdleTimeout,
		e.node.WriteTimeout,
		nconn)

	return label, conn, nil
}
