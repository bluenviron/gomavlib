package gomavlib

import (
	"fmt"
	"io"
	"net"
)

type transportServerConf interface {
	GetAddress() string
}

// TransportTcpServer sends and reads frames through a TCP server.
// TCP can be good for routing frames through the internet, but it is not the most
// appropriate way for transfering frames from a UAV to a GCS, since it does
// not allow frame losses.
type TransportTcpServer struct {
	// listen address, i.e 0.0.0.0:5600
	Address string
}

func (c TransportTcpServer) GetAddress() string {
	return c.Address
}

// TransportUdpServer sends and reads frames through an UDP server.
// This is the most appropriate way for transfering frames from a UAV to a GPS
// if they are connected to the same network.
type TransportUdpServer struct {
	// listen address, i.e 0.0.0.0:5600
	Address string
}

func (c TransportUdpServer) GetAddress() string {
	return c.Address
}

type transportServer struct {
	conf      transportServerConf
	node      *Node
	terminate chan struct{}
	isUdp     bool
	listener  net.Listener
}

func (conf TransportTcpServer) init(n *Node) (transport, error) {
	return initTransportServer(n, false, conf)
}

func (conf TransportUdpServer) init(n *Node) (transport, error) {
	return initTransportServer(n, true, conf)
}

func initTransportServer(node *Node, isUdp bool, conf transportServerConf) (transport, error) {
	_, _, err := net.SplitHostPort(conf.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	var listener net.Listener
	if isUdp == true {
		listener, err = newUdpListener("udp4", conf.GetAddress())
	} else {
		listener, err = net.Listen("tcp4", conf.GetAddress())
	}
	if err != nil {
		return nil, err
	}

	t := &transportServer{
		conf:      conf,
		node:      node,
		terminate: make(chan struct{}),
		isUdp:     isUdp,
		listener:  listener,
	}
	return t, nil
}

func (t *transportServer) closePrematurely() {
	t.listener.Close()
}

func (t *transportServer) do() {
	listenDone := make(chan struct{})

	go func() {
		defer func() { listenDone <- struct{}{} }()

		for {
			rawConn, err := t.listener.Accept()
			if err != nil {
				break
			}

			conn := &netTimedConn{rawConn}
			tconn := &TransportChannel{
				transport: t,
				writer:    conn,
			}
			t.node.addTransportChannel(tconn)

			t.node.wg.Add(1)
			go func() {
				defer t.node.wg.Done()
				defer t.node.removeTransportChannel(tconn)
				defer conn.Close()

				var buf [netBufferSize]byte
				for {
					n, err := conn.Read(buf[:])
					if err != nil {
						break
					}
					t.node.processBuffer(tconn, buf[:n])
				}
			}()
		}
	}()

	select {
	// unexpected error, wait for terminate
	case <-listenDone:
		t.listener.Close()
		<-t.node.terminate

	// terminated
	case <-t.node.terminate:
		t.listener.Close()
		<-listenDone
	}

	// close all channels
	func() {
		t.node.mutex.Lock()
		defer t.node.mutex.Unlock()

		for conn := range t.node.channels {
			if conn.transport == t {
				conn.writer.(io.Closer).Close()
			}
		}
	}()
}
