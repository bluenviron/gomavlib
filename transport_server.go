package gomavlib

import (
	"fmt"
	"io"
	"net"
)

type transportServerConf interface {
	isUdp() bool
	getAddress() string
}

// TransportTcpServer reads and writes frames through a TCP server.
// TCP can be good for routing frames through the internet, but it is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type TransportTcpServer struct {
	// listen address, i.e 0.0.0.0:5600
	Address string
}

func (TransportTcpServer) isUdp() bool {
	return false
}

func (c TransportTcpServer) getAddress() string {
	return c.Address
}

// TransportUdpServer reads and writes frames through an UDP server.
// This is the most appropriate way for transferring frames from a UAV to a GPS
// if they are connected to the same network.
type TransportUdpServer struct {
	// listen address, i.e 0.0.0.0:5600
	Address string
}

func (TransportUdpServer) isUdp() bool {
	return true
}

func (c TransportUdpServer) getAddress() string {
	return c.Address
}

type transportServer struct {
	conf      transportServerConf
	listener  net.Listener
	terminate chan struct{}
}

func (conf TransportTcpServer) init() (transport, error) {
	return initTransportServer(conf)
}

func (conf TransportUdpServer) init() (transport, error) {
	return initTransportServer(conf)
}

func initTransportServer(conf transportServerConf) (transport, error) {
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

	t := &transportServer{
		conf:      conf,
		listener:  listener,
		terminate: make(chan struct{}, 1),
	}
	return t, nil
}

func (t *transportServer) isTransport() {
}

func (t *transportServer) Close() error {
	t.terminate <- struct{}{}
	t.listener.Close()
	return nil
}

func (t *transportServer) Accept() (io.ReadWriteCloser, error) {
	rawConn, err := t.listener.Accept()

	// wait termination, do not report errors
	if err != nil {
		<-t.terminate
		return nil, errorTerminated
	}

	conn := &netTimedConn{rawConn}

	return conn, nil
}
