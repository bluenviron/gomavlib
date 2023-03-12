package gomavlib

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aler9/gomavlib/pkg/timednetconn"
	"github.com/aler9/gomavlib/pkg/udplistener"
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

type endpointServer struct {
	conf         endpointServerConf
	listener     net.Listener
	writeTimeout time.Duration

	// in
	terminate chan struct{}
}

func (conf EndpointTCPServer) init(node *Node) (Endpoint, error) {
	return initEndpointServer(node, conf)
}

func (conf EndpointUDPServer) init(node *Node) (Endpoint, error) {
	return initEndpointServer(node, conf)
}

func initEndpointServer(node *Node, conf endpointServerConf) (Endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	var listener net.Listener
	if conf.isUDP() {
		listener, err = udplistener.New("udp4", conf.getAddress())
	} else {
		listener, err = net.Listen("tcp4", conf.getAddress())
	}
	if err != nil {
		return nil, err
	}

	t := &endpointServer{
		conf:         conf,
		writeTimeout: node.conf.WriteTimeout,
		listener:     listener,
		terminate:    make(chan struct{}),
	}
	return t, nil
}

func (t *endpointServer) isEndpoint() {}

func (t *endpointServer) Conf() EndpointConf {
	return t.conf
}

func (t *endpointServer) close() error {
	close(t.terminate)
	t.listener.Close()
	return nil
}

func (t *endpointServer) accept() (string, io.ReadWriteCloser, error) {
	nconn, err := t.listener.Accept()
	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return "", nil, errTerminated
	}

	label := fmt.Sprintf("%s:%s", func() string {
		if t.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), nconn.RemoteAddr())

	conn := timednetconn.New(t.writeTimeout, nconn)

	return label, conn, nil
}
