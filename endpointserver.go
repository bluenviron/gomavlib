package gomavlib

import (
	"fmt"
	"io"
	"net"

	"github.com/aler9/gomavlib/pkg/udplistener"
)

type endpointServerConf interface {
	isUDP() bool
	getAddress() string
	init() (Endpoint, error)
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
	conf     endpointServerConf
	listener net.Listener

	// in
	terminate chan struct{}
}

func (conf EndpointTCPServer) init() (Endpoint, error) {
	return initEndpointServer(conf)
}

func (conf EndpointUDPServer) init() (Endpoint, error) {
	return initEndpointServer(conf)
}

func initEndpointServer(conf endpointServerConf) (Endpoint, error) {
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
		conf:      conf,
		listener:  listener,
		terminate: make(chan struct{}),
	}
	return t, nil
}

func (t *endpointServer) isEndpoint() {}

func (t *endpointServer) Conf() EndpointConf {
	return t.conf
}

func (t *endpointServer) Close() error {
	close(t.terminate)
	t.listener.Close()
	return nil
}

func (t *endpointServer) Accept() (string, io.ReadWriteCloser, error) {
	rawConn, err := t.listener.Accept()
	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return "", nil, errorTerminated
	}

	label := fmt.Sprintf("%s:%s", func() string {
		if t.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), rawConn.RemoteAddr())

	conn := &netTimedConn{rawConn}

	return label, conn, nil
}
