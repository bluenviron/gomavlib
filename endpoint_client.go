package gomavlib

import (
	"context"
	"io"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

var reconnectPeriod = 2 * time.Second

type endpointClient struct {
	node *Node
	conf EndpointCustomClient

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
	timedContext, timedContextClose := context.WithTimeout(e.ctx, e.node.ReadTimeout)
	nconn, err := e.conf.Connect(timedContext)
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
	if e.conf.Label != "" {
		return e.conf.Label
	}
	return "custom"
}
