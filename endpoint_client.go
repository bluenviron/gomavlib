package gomavlib

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

var reconnectPeriod = 2 * time.Second

type endpointClientConf interface {
	init(*Node) (Endpoint, error)
}

// EndpointTCPClient sets up a endpoint that works with a TCP client.
// TCP is fit for routing frames through the internet, but is not the most
// appropriate way for transferring frames from a UAV to a GCS, since it does
// not allow frame losses.
type EndpointTCPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (conf EndpointTCPClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

// EndpointUDPClient sets up a endpoint that works with a UDP client.
type EndpointUDPClient struct {
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
}

func (conf EndpointUDPClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

// EndpointCustomClient sets up a endpoint that works with a custom implementation
// by providing a Connect func that returns a net.Conn.
type EndpointCustomClient struct {
	// custom connect function that opens the connection
	Connect func(ctx context.Context) (net.Conn, error)

	// the label of the protocol
	Label string
}

func (conf EndpointCustomClient) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

type endpointClient struct {
	node *Node
	conf endpointClientConf

	ctx       context.Context
	ctxCancel func()
	first     bool
}

func (e *endpointClient) initialize() error {
	e.ctx, e.ctxCancel = context.WithCancel(context.Background())

	return nil
}

func (e *endpointClient) isEndpoint() {}

func (e *endpointClient) Conf() EndpointConf {
	return e.conf
}

func (e *endpointClient) close() {
	e.ctxCancel()
}

func (e *endpointClient) oneChannelAtAtime() bool {
	return true
}

func (e *endpointClient) connect() (io.ReadWriteCloser, error) {
	// in UDP, the only possible error is a DNS failure
	// in TCP, the handshake must be completed
	timedContext, timedContextClose := context.WithTimeout(e.ctx, e.node.ReadTimeout)
	nconn, err := func() (net.Conn, error) {
		switch conf := e.conf.(type) {
		case EndpointTCPClient:
			return (&net.Dialer{}).DialContext(timedContext, "tcp4", conf.Address)

		case EndpointUDPClient:
			return (&net.Dialer{}).DialContext(timedContext, "udp4", conf.Address)

		case EndpointCustomClient:
			return conf.Connect(timedContext)

		default:
			panic("should not happen")
		}
	}()
	timedContextClose()

	if err != nil {
		return nil, err
	}

	return timednetconn.New(
		e.node.IdleTimeout,
		e.node.WriteTimeout,
		nconn,
	), nil
}

func (e *endpointClient) provide() (string, io.ReadWriteCloser, error) {
	if !e.first {
		e.first = true
	} else {
		select {
		case <-time.After(reconnectPeriod):
		case <-e.ctx.Done():
			return "", nil, errTerminated
		}
	}

	for {
		conn, err := e.connect()
		if err != nil {
			select {
			case <-time.After(reconnectPeriod):
				continue
			case <-e.ctx.Done():
				return "", nil, errTerminated
			}
		}

		return e.label(), conn, nil
	}
}

func (e *endpointClient) label() string {
	switch conf := e.conf.(type) {
	case EndpointTCPClient:
		return "tcp:" + conf.Address
	case EndpointUDPClient:
		return "udp:" + conf.Address
	case EndpointCustomClient:
		if conf.Label != "" {
			return conf.Label
		}
		return "custom"
	default:
		return "unknown"
	}
}
