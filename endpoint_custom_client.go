package gomavlib

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/bluenviron/gomavlib/v4/pkg/timednetconn"
)

var _ Endpoint = (*EndpointCustomClient)(nil)

var reconnectPeriod = 2 * time.Second

// EndpointCustomClient is an endpoint that works with a custom implementation
// by providing a Connect func that returns a net.Conn.
type EndpointCustomClient struct {
	// custom connect function that opens the connection
	Connect func(ctx context.Context) (net.Conn, error)

	// the label of the protocol
	Label string

	// whether the connection is datagram-based (e.g. UDP).
	IsDatagram bool

	node      *Node
	ctx       context.Context
	ctxCancel func()
	first     bool
}

func (e *EndpointCustomClient) init(node *Node) error {
	e.node = node
	e.ctx, e.ctxCancel = context.WithCancel(context.Background())
	return nil
}

func (e *EndpointCustomClient) isEndpoint() {}

func (e *EndpointCustomClient) close() {
	e.ctxCancel()
}

func (e *EndpointCustomClient) oneChannelAtAtime() bool {
	return true
}

func (e *EndpointCustomClient) isDatagram() bool {
	return e.IsDatagram
}

func (e *EndpointCustomClient) connect() (io.ReadWriteCloser, error) {
	timedContext, timedContextClose := context.WithTimeout(e.ctx, e.node.ReadTimeout)
	nconn, err := e.Connect(timedContext)
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

func (e *EndpointCustomClient) provide() (string, io.ReadWriteCloser, error) {
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

func (e *EndpointCustomClient) label() string {
	if e.Label != "" {
		return e.Label
	}
	return "custom"
}
