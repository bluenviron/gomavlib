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
// TCP can be good for routing frames through the internet, but it is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTcpServer struct {
	// listen address, i.e 0.0.0.0:5600
	Address string
}

func (EndpointTcpServer) isUdp() bool {
	return false
}

func (conf EndpointTcpServer) getAddress() string {
	return conf.Address
}

// EndpointUdpServer sets up a endpoint that works through an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GPS
// if they are connected to the same network.
type EndpointUdpServer struct {
	// listen address, i.e 0.0.0.0:5600
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

func (conf EndpointTcpServer) init() (endpoint, error) {
	return initEndpointServer(conf)
}

func (conf EndpointUdpServer) init() (endpoint, error) {
	return initEndpointServer(conf)
}

func initEndpointServer(conf endpointServerConf) (endpoint, error) {
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
		terminate: make(chan struct{}, 1),
	}
	return t, nil
}

func (t *endpointServer) isEndpoint() {
}

func (t *endpointServer) Close() error {
	t.terminate <- struct{}{}
	t.listener.Close()
	return nil
}

func (t *endpointServer) Accept() (io.ReadWriteCloser, error) {
	rawConn, err := t.listener.Accept()

	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return nil, errorTerminated
	}

	conn := &netTimedConn{rawConn}

	return conn, nil
}
