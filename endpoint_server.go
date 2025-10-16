package gomavlib

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
	"github.com/pion/transport/v2/udp"
)

type endpointServerType int

const (
	endpointServerTypeTCP endpointServerType = iota
	endpointServerTypeUDP
	endpointServerTypeCustom
)

type endpointServerConf interface {
	serverType() endpointServerType
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

func (EndpointTCPServer) serverType() endpointServerType {
	return endpointServerTypeTCP
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

func (EndpointUDPServer) serverType() endpointServerType {
	return endpointServerTypeUDP
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

func (EndpointCustomServer) serverType() endpointServerType {
	return endpointServerTypeCustom
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

	switch e.conf.serverType() {
	case endpointServerTypeUDP:
		var addr *net.UDPAddr
		addr, err = net.ResolveUDPAddr("udp4", e.conf.getAddress())
		if err != nil {
			return err
		}

		e.listener, err = udp.Listen("udp4", addr)
		if err != nil {
			return err
		}
	case endpointServerTypeTCP:
		e.listener, err = net.Listen("tcp4", e.conf.getAddress())
		if err != nil {
			return err
		}
	case endpointServerTypeCustom:
		if customConf, ok := e.conf.(EndpointCustomServer); ok {
			e.listener, err = customConf.Listen(e.conf.getAddress())
			if err != nil {
				return err
			}
		} else {
			return errors.New("type assertion error to endpointcustomserver")
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
		switch e.conf.serverType() {
		case endpointServerTypeTCP:
			return "tcp"
		case endpointServerTypeUDP:
			return "udp"
		case endpointServerTypeCustom:
			if customConf, ok := e.conf.(EndpointCustomServer); ok {
				if customConf.Label != "" {
					return customConf.Label
				}
			}
			return "cust"
		default:
			return "unk"
		}
	}(), nconn.RemoteAddr())

	conn := timednetconn.New(
		e.node.IdleTimeout,
		e.node.WriteTimeout,
		nconn)

	return label, conn, nil
}
