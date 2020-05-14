package gomavlib

import (
	"fmt"
	"io"
	"net"
)

type endpointServerConf interface {
	isUdp() bool
	getAddress() string
}

// EndpointTcpServer sets up a endpoint that works through a TCP server.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTcpServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointTcpServer) isUdp() bool {
	return false
}

func (conf EndpointTcpServer) getAddress() string {
	return conf.Address
}

// EndpointUdpServer sets up a endpoint that works through an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GCS
// if they are connected to the same network.
type EndpointUdpServer struct {
	// listen address, example: 0.0.0.0:5600
	Address string
}

func (EndpointUdpServer) isUdp() bool {
	return true
}

func (conf EndpointUdpServer) getAddress() string {
	return conf.Address
}

type endpointServer struct {
	conf      endpointServerConf
	listener  net.Listener
	terminate chan struct{}
}

func (conf EndpointTcpServer) init() (Endpoint, error) {
	return initEndpointServer(conf)
}

func (conf EndpointUdpServer) init() (Endpoint, error) {
	return initEndpointServer(conf)
}

func initEndpointServer(conf endpointServerConf) (Endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	var listener net.Listener
	if conf.isUdp() == true {
		listener, err = newUdpListener("udp4", conf.getAddress())
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

func (t *endpointServer) Conf() interface{} {
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
		if t.conf.isUdp() {
			return "udp"
		}
		return "tcp"
	}(), rawConn.RemoteAddr())

	conn := &netTimedConn{rawConn}

	return label, conn, nil
}
