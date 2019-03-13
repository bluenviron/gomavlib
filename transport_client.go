package gomavlib

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type transportClientConf interface {
	isUdp() bool
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
	conf       transportClientConf
	mutex      sync.Mutex
	terminated bool
	terminate  chan struct{}
	conn       io.ReadWriteCloser
}

func initTransportClient(node *Node, conf transportClientConf) (transport, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	t := &transportClient{
		conf:      conf,
		terminate: make(chan struct{}, 1),
	}

	tc := &TransportCustom{
		ReadWriteCloser: t,
	}
	return tc.init(node)
}

func (t *transportClient) Close() error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.terminated == true {
		return nil
	}
	t.terminate <- struct{}{}
	return nil
}

func (t *transportClient) Write(buf []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.conn != nil {
		return t.conn.Write(buf)
	}
	return 0, fmt.Errorf("disconnected")
}

func (t *transportClient) Read(buf []byte) (int, error) {
	for {
		// mutex is not necessary since Read() is the only one that can fill and empty conn
		if t.conn == nil {
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
			case <-t.terminate:
				return 0, fmt.Errorf("terminated")
			}

			// wait some seconds before reconnecting
			if rawConn == nil {
				timer := time.NewTimer(netReconnectPeriod)
				select {
				case <-timer.C:
					continue
				case <-t.terminate:
					return 0, fmt.Errorf("terminated")
				}
			}

			func() {
				t.mutex.Lock()
				defer t.mutex.Unlock()
				t.conn = &netTimedConn{rawConn}
			}()
		}

		var n int
		var err error
		readDone := make(chan struct{})
		go func() {
			defer func() { readDone <- struct{}{} }()
			n, err = t.conn.Read(buf)
		}()

		select {
		case <-readDone:
		case <-t.terminate:
			t.conn.Close()
			<-readDone
			return 0, fmt.Errorf("terminated")
		}

		// unexpected error, restart connection
		if err != nil {
			t.conn.Close()
			func() {
				t.mutex.Lock()
				defer t.mutex.Unlock()
				t.conn = nil
			}()
			continue
		}

		return n, nil
	}
}
