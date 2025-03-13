package gomavlib

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/bluenviron/gomavlib/v3/pkg/reconnector"
	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

type endpointClientConf interface {
	isUDP() bool
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

func (EndpointTCPClient) isUDP() bool {
	return false
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

func (EndpointUDPClient) isUDP() bool {
	return true
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

type endpointClient struct {
	node *Node
	conf endpointClientConf

	reconnector *reconnector.Reconnector
}

func (e *endpointClient) initialize() error {
	_, _, err := net.SplitHostPort(e.conf.getAddress())
	if err != nil {
		return fmt.Errorf("invalid address")
	}

	e.reconnector = reconnector.New(
		func(ctx context.Context) (io.ReadWriteCloser, error) {
			network := func() string {
				if e.conf.isUDP() {
					return "udp4"
				}
				return "tcp4"
			}()

			// in UDP, the only possible error is a DNS failure
			// in TCP, the handshake must be completed
			timedContext, timedContextClose := context.WithTimeout(ctx, e.node.ReadTimeout)
			nconn, err := (&net.Dialer{}).DialContext(timedContext, network, e.conf.getAddress())
			timedContextClose()

			if err != nil {
				return nil, err
			}

			return timednetconn.New(
				e.node.IdleTimeout,
				e.node.WriteTimeout,
				nconn,
			), nil
		},
	)

	return nil
}

func (e *endpointClient) isEndpoint() {}

func (e *endpointClient) Conf() EndpointConf {
	return e.conf
}

func (e *endpointClient) close() {
	e.reconnector.Close()
}

func (e *endpointClient) oneChannelAtAtime() bool {
	return true
}

func (e *endpointClient) provide() (string, io.ReadWriteCloser, error) {
	conn, ok := e.reconnector.Reconnect()
	if !ok {
		return "", nil, errTerminated
	}

	return e.label(), conn, nil
}

func (e *endpointClient) label() string {
	return fmt.Sprintf("%s:%s", func() string {
		if e.conf.isUDP() {
			return "udp"
		}
		return "tcp"
	}(), e.conf.getAddress())
}
