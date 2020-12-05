package gomavlib

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/aler9/gomavlib/pkg/multibuffer"
)

type endpointClientConf interface {
	isUDP() bool
	getAddress() string
	init() (Endpoint, error)
}

// EndpointTCPClient sets up a endpoint that works with a TCP client.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (EndpointTCPClient) isUDP() bool {
	return false
}

func (conf EndpointTCPClient) getAddress() string {
	return conf.Address
}

func (conf EndpointTCPClient) init() (Endpoint, error) {
	return initEndpointClient(conf)
}

// EndpointUDPClient sets up a endpoint that works with a UDP client.
type EndpointUDPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (EndpointUDPClient) isUDP() bool {
	return true
}

func (conf EndpointUDPClient) getAddress() string {
	return conf.Address
}

func (conf EndpointUDPClient) init() (Endpoint, error) {
	return initEndpointClient(conf)
}

type endpointClient struct {
	conf        endpointClientConf
	writerMutex sync.Mutex
	writer      io.Writer

	// in
	terminate chan struct{}
	read      chan []byte
}

func initEndpointClient(conf endpointClientConf) (Endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	t := &endpointClient{
		conf:      conf,
		terminate: make(chan struct{}),
		read:      make(chan []byte),
	}

	// work in a separate routine
	// in this way we connect immediately, not after the first Read()
	go t.do()
	return t, nil
}

func (t *endpointClient) isEndpoint() {}

func (t *endpointClient) Conf() EndpointConf {
	return t.conf
}

func (t *endpointClient) Label() string {
	return fmt.Sprintf("%s:%s", func() string {
		if t.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), t.conf.getAddress())
}

func (t *endpointClient) Close() error {
	close(t.terminate)
	return nil
}

func (t *endpointClient) do() {
	mb := multibuffer.New(2, bufferSize)

	for {
		// solve address and connect
		// in UDP, the only possible error is a DNS failure
		// in TCP, the handshake must be completed
		var rawConn net.Conn
		dialDone := make(chan struct{}, 1)
		go func() {
			defer close(dialDone)

			network := func() string {
				if t.conf.isUDP() {
					return "udp4"
				}
				return "tcp4"
			}()

			var err error
			rawConn, err = net.DialTimeout(network, t.conf.getAddress(), netConnectTimeout)
			if err != nil {
				rawConn = nil // ensure rawConn is nil in case of error
			}
		}()

		select {
		case <-dialDone:
		case <-t.terminate:
			go func() {
				for range t.read {
				}
			}()
			close(t.read)
			return
		}

		if rawConn == nil {
			ok := func() bool {
				// wait some seconds before reconnecting
				timer := time.NewTimer(netReconnectPeriod)
				defer timer.Stop()

				select {
				case <-timer.C:
					return true
				case <-t.terminate:
					go func() {
						for range t.read {
						}
					}()
					close(t.read)
					return false
				}
			}()
			if !ok {
				return
			}
			continue
		}

		conn := &netTimedConn{rawConn}
		func() {
			t.writerMutex.Lock()
			defer t.writerMutex.Unlock()
			t.writer = conn
		}()

		readerDone := make(chan struct{})
		go func() {
			defer close(readerDone)

			for {
				buf := mb.Next()
				n, err := conn.Read(buf)
				if err != nil {
					return
				}

				t.read <- buf[:n]
			}
		}()

		select {
		case <-readerDone:
		case <-t.terminate:
			go func() {
				for range t.read {
				}
			}()
			conn.Close()
			<-readerDone
			close(t.read)
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
	src, ok := <-t.read
	if !ok {
		return 0, errorTerminated
	}
	n := copy(buf, src)
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
