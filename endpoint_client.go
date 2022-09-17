package gomavlib

import (
	"context"
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
	conf endpointClientConf

	ctx         context.Context
	ctxCancel   func()
	mb          *multibuffer.MultiBuffer
	writerMutex sync.Mutex
	writer      io.Writer

	// in
	read chan []byte
}

func initEndpointClient(conf endpointClientConf) (Endpoint, error) {
	_, _, err := net.SplitHostPort(conf.getAddress())
	if err != nil {
		return nil, fmt.Errorf("invalid address")
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	t := &endpointClient{
		conf:      conf,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		mb:        multibuffer.New(2, bufferSize),
		read:      make(chan []byte),
	}

	go t.run()

	return t, nil
}

func (t *endpointClient) isEndpoint() {}

func (t *endpointClient) Conf() EndpointConf {
	return t.conf
}

func (t *endpointClient) label() string {
	return fmt.Sprintf("%s:%s", func() string {
		if t.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), t.conf.getAddress())
}

func (t *endpointClient) Close() error {
	t.ctxCancel()
	return nil
}

func (t *endpointClient) run() {
	for {
		t.runInner()

		select {
		case <-time.After(netReconnectPeriod):
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *endpointClient) runInner() error {
	// solve address and connect
	// in UDP, the only possible error is a DNS failure
	// in TCP, the handshake must be completed
	network := func() string {
		if t.conf.isUDP() {
			return "udp4"
		}
		return "tcp4"
	}()
	timedContext, timedContextClose := context.WithTimeout(t.ctx, netConnectTimeout)
	nconn, err := (&net.Dialer{}).DialContext(timedContext, network, t.conf.getAddress())
	timedContextClose()
	if err != nil {
		return err
	}

	conn := &netTimedConn{nconn}
	func() {
		t.writerMutex.Lock()
		defer t.writerMutex.Unlock()
		t.writer = conn
	}()

	readerDone := make(chan error)
	go func() {
		readerDone <- func() error {
			for {
				buf := t.mb.Next()
				n, err := conn.Read(buf)
				if err != nil {
					return err
				}

				select {
				case t.read <- buf[:n]:
				case <-t.ctx.Done():
					return errTerminated
				}
			}
		}()
	}()

	select {
	case err := <-readerDone: // unexpected error
		conn.Close()
		func() {
			t.writerMutex.Lock()
			defer t.writerMutex.Unlock()
			t.writer = nil
		}()
		return err

	case <-t.ctx.Done(): // Close() has been called
		conn.Close()
		<-readerDone
		return errTerminated
	}
}

func (t *endpointClient) Read(buf []byte) (int, error) {
	select {
	case src := <-t.read:
		n := copy(buf, src)
		return n, nil

	case <-t.ctx.Done():
		return 0, errTerminated
	}
}

func (t *endpointClient) Write(buf []byte) (int, error) {
	t.writerMutex.Lock()
	defer t.writerMutex.Unlock()

	if t.writer == nil {
		return 0, fmt.Errorf("disconnected")
	}

	return t.writer.Write(buf)
}
