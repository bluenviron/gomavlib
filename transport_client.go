package gomavlib

import (
	"fmt"
	"net"
	"time"
)

type transportClientConf interface {
	getAddress() string
}

// TransportTcpClient reads and writes frames through a TCP client.
type TransportTcpClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (TransportTcpClient) isUdp() bool {
	return false
}

func (conf TransportTcpClient) getAddress() string {
	return conf.Address
}

func (conf TransportTcpClient) init(n *Node) (transport, error) {
	return initTransportClient(n, conf)
}

// TransportUdpClient reads and writes frames through a UDP client.
type TransportUdpClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (TransportUdpClient) isUdp() bool {
	return true
}

func (conf TransportUdpClient) getAddress() string {
	return conf.Address
}

func (conf TransportUdpClient) init(n *Node) (transport, error) {
	return initTransportClient(n, conf)
}

type transportClient struct {
	conf      transportClientConf
	node      *Node
	terminate chan struct{}
}

func initTransportClient(node *Node, conf transportClientConf) (transport, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	t := &transportClient{
		conf:  conf,
		node:  node,
	}
	return t, nil
}

func (t *transportClient) closePrematurely() {
}

func (t *transportClient) do() {
	// reuse the same tconn such that it can be used as map key
	tconn := &TransportChannel{
		transport: t,
		writer:    nil,
	}

	for {
		// solve address and connect
		// in UDP, the only possible error is a DNS failure
		// in TCP, the handshake must be completed
		var rawConn net.Conn
		dialDone := make(chan struct{}, 1)
		go func() {
			var network string
			if t.conf.isUdp() == true {
				network = "udp4"
			} else {
				network = "tcp4"
			}
			var err error
			rawConn, err = net.DialTimeout(network, t.conf.getAddress(), netConnectTimeout)
			if err != nil {
				rawConn = nil // ensure rawConn is nil in case of error
			}
			dialDone <- struct{}{}
		}()

		select {
		case <-dialDone:
		case <-t.node.terminate:
			return
		}

		// wait some seconds before reconnecting
		if rawConn == nil {
			timer := time.NewTimer(netReconnectPeriod)
			select {
			case <-timer.C:
				continue
			case <-t.node.terminate:
				return
			}
		}

		conn := &netTimedConn{rawConn}
		tconn.writer = conn
		t.node.addTransportChannel(tconn)

		readDone := make(chan struct{})
		go func() {
			defer func() { readDone <- struct{}{} }()
			defer t.node.removeTransportChannel(tconn)

			var buf [netBufferSize]byte
			for {
				n, err := conn.Read(buf[:])
				if err != nil {
					break
				}
				t.node.processBuffer(tconn, buf[:n])
			}
		}()

		select {
		// unexpected error, restart connection
		case <-readDone:
			conn.Close()
			continue

		// terminated
		case <-t.node.terminate:
			conn.Close()
			<-readDone
			return
		}
	}
}
