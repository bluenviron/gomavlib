package gomavlib

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type endpointClientConf interface {
	isUdp() bool
	getAddress() string
}

// EndpointTcpClient sets up a endpoint that works through a TCP client.
type EndpointTcpClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (EndpointTcpClient) isUdp() bool {
	return false
}

func (conf EndpointTcpClient) getAddress() string {
	return conf.Address
}

func (conf EndpointTcpClient) init() (endpoint, error) {
	return initEndpointClient(conf)
}

// EndpointUdpClient sets up a endpoint that works through a UDP client.
type EndpointUdpClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (EndpointUdpClient) isUdp() bool {
	return true
}

func (conf EndpointUdpClient) getAddress() string {
	return conf.Address
}

func (conf EndpointUdpClient) init() (endpoint, error) {
	return initEndpointClient(conf)
}

type endpointClient struct {
	conf      endpointClientConf
	mutex     sync.Mutex
	terminate chan struct{}
	conn      io.ReadWriteCloser
}

func initEndpointClient(conf endpointClientConf) (endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	t := &endpointClient{
		conf:      conf,
		terminate: make(chan struct{}, 1),
	}
	return t, nil
}

func (*endpointClient) isEndpoint() {}

func (t *endpointClient) Close() error {
	t.terminate <- struct{}{}
	return nil
}

func (t *endpointClient) Write(buf []byte) (int, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.conn != nil {
		return t.conn.Write(buf)
	}
	return 0, fmt.Errorf("disconnected")
}

func (t *endpointClient) Read(buf []byte) (int, error) {
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
				return 0, errorTerminated
			}

			// wait some seconds before reconnecting
			if rawConn == nil {
				timer := time.NewTimer(netReconnectPeriod)
				select {
				case <-timer.C:
					continue
				case <-t.terminate:
					return 0, errorTerminated
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
			return 0, errorTerminated
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
