package gomavlib

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
	"github.com/pion/transport/v2/udp"
)

type endpointServerConf interface {
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

// EndpointCustomServer sets up a endpoint that works with custom implementations
// by providing a custom Listen func that returns a net.Listener.
// This allows you to use custom protocols that conform to the net.listner.
// A use case could be to add encrypted protocol implementations like DTLS or TCP with TLS.
type EndpointCustomServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
	// function to invoke when server should start listening
	Listen func(address string) (net.Listener, error)
	// the label of the protocol
	Label string
}

func (conf EndpointCustomServer) getAddress() string {
	return conf.Address
}

func (conf EndpointCustomServer) init(node *Node) (Endpoint, error) {
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

	switch conf := e.conf.(type) {
	case EndpointUDPServer:
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp4", e.conf.getAddress())
		if err != nil {
			return err
		}

		e.listener, err = udp.Listen("udp4", addr)
		if err != nil {
			return err
		}

	case EndpointTCPServer:
		e.listener, err = net.Listen("tcp4", e.conf.getAddress())
		if err != nil {
			return err
		}

	case EndpointCustomServer:
		e.listener, err = conf.Listen(e.conf.getAddress())
		if err != nil {
			return err
		}

	default:
		return errors.New("unsupported server-type")
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
		switch conf := e.conf.(type) {
		case EndpointTCPServer:
			return "tcp"
		case EndpointUDPServer:
			return "udp"
		case EndpointCustomServer:
			if conf.Label != "" {
				return conf.Label
			}
			return "custom"
		default:
			return "unknown"
		}
	}(), nconn.RemoteAddr())

	conn := timednetconn.New(
		e.node.IdleTimeout,
		e.node.WriteTimeout,
		nconn)

	return label, conn, nil
}
