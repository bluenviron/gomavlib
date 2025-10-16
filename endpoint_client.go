package gomavlib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

var reconnectPeriod = 2 * time.Second

type endpointClientConf interface {
	clientType() endpointServerType
	getAddress() string
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

func (EndpointTCPClient) clientType() endpointServerType {
	return endpointServerTypeTCP
}

func (conf EndpointTCPClient) getAddress() string {
	return conf.Address
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

func (EndpointUDPClient) clientType() endpointServerType {
	return endpointServerTypeUDP
}

func (conf EndpointUDPClient) getAddress() string {
	return conf.Address
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
	// domain name or IP of the server to connect to, example: 1.2.3.4:5600
	Address string
	// custom connect function that connects to the provided address
	Connect func(address string) (net.Conn, error)
	// the label of the protocol
	Label string
}

func (EndpointCustomClient) clientType() endpointServerType {
	return endpointServerTypeCustom
}

func (conf EndpointCustomClient) getAddress() string {
	return conf.Address
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
	_, _, err := net.SplitHostPort(e.conf.getAddress())
	if err != nil {
		return fmt.Errorf("invalid address")
	}

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
	network := func() string {
		switch e.conf.clientType() {
		case endpointServerTypeTCP:
			return "tcp4"
		case endpointServerTypeUDP:
			return "udp4"
		case endpointServerTypeCustom:
			return "cust"
		default:
			return ""
		}
	}()

	// in UDP, the only possible error is a DNS failure
	// in TCP, the handshake must be completed
	timedContext, timedContextClose := context.WithTimeout(e.ctx, e.node.ReadTimeout)
	nconn, err := func() (net.Conn, error) {
		if network == "cust" {
			if customConf, ok := e.conf.(EndpointCustomClient); ok {
				if customConf.Connect == nil {
					return nil, errors.New("no connect function provided on custom endpoint")
				}
				return customConf.Connect(customConf.Address)
			}
			return nil, errors.New("failed type assertion to endpointcustomclient")
		}
		return (&net.Dialer{}).DialContext(timedContext, network, e.conf.getAddress())
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
	return fmt.Sprintf("%s:%s", func() string {
		switch e.conf.clientType() {
		case endpointServerTypeTCP:
			return "tcp"
		case endpointServerTypeUDP:
			return "udp"
		case endpointServerTypeCustom:
			if customConf, ok := e.conf.(EndpointCustomClient); ok {
				if customConf.Label != "" {
					return customConf.Label
				}
			}
			return "cust"
		default:
			return "unk"
		}
	}(), e.conf.getAddress())
}
