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
	conf        endpointClientConf
	terminate   chan struct{}
	readChan    chan []byte
	readDone    chan struct{}
	writerMutex sync.Mutex
	writer      io.Writer
}

func initEndpointClient(conf endpointClientConf) (endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	t := &endpointClient{
		conf:      conf,
		terminate: make(chan struct{}, 1),
		readChan:  make(chan []byte),
		readDone:  make(chan struct{}),
	}

	// work in a separate routine
	// in this way we connect immediately, not after the first Read()
	go t.do()
	return t, nil
}

func (*endpointClient) isEndpoint() {}

func (t *endpointClient) Close() error {
	t.terminate <- struct{}{}
	return nil
}

func (t *endpointClient) do() {
	defer func() {
		// consumate readChan in case user is not calling Read()
		go func() {
			for range t.readChan {
			}
		}()
		close(t.readChan)
	}()

	buf := make([]byte, netBufferSize)

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
		case <-t.terminate:
			return
		}

		// wait some seconds before reconnecting
		if rawConn == nil {
			timer := time.NewTimer(netReconnectPeriod)
			select {
			case <-timer.C:
				continue
			case <-t.terminate:
				return
			}
		}

		conn := &netTimedConn{rawConn}
		func() {
			t.writerMutex.Lock()
			defer t.writerMutex.Unlock()
			t.writer = conn
		}()

		var n int
		var err error
		cycleDone := make(chan struct{})
		go func() {
			defer func() { cycleDone <- struct{}{} }()
			for {
				n, err = conn.Read(buf)
				if err != nil {
					return
				}
				t.readChan <- buf[:n]
				<-t.readDone
			}
		}()

		select {
		case <-cycleDone:
		case <-t.terminate:
			conn.Close()
			<-cycleDone
			return
		}

		// unexpected error, restart connection
		conn.Close()
		func() {
			t.writerMutex.Lock()
			defer t.writerMutex.Unlock()
			t.writer = nil
		}()
	}
}

func (t *endpointClient) Read(buf []byte) (int, error) {
	src, ok := <-t.readChan
	if ok == false {
		return 0, errorTerminated
	}
	n := copy(buf, src)
	t.readDone <- struct{}{}
	return n, nil
}

func (t *endpointClient) Write(buf []byte) (int, error) {
	t.writerMutex.Lock()
	defer t.writerMutex.Unlock()

	// drop packets if disconnected
	if t.writer == nil {
		return 0, fmt.Errorf("disconnected")
	}

	return t.writer.Write(buf)
}
