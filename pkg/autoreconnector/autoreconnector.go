// Package autoreconnector contains a io.ReadWriteCloser wrapper that implements automatic reconnection.
package autoreconnector

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

var (
	reconnectPeriod = 2 * time.Second
	errTerminated   = errors.New("terminated")
)

type autoreconnector struct {
	connect func(context.Context) (io.ReadWriteCloser, error)

	ctx       context.Context
	ctxCancel func()
	conn      io.ReadWriteCloser
	connMutex sync.Mutex
}

// New returns a io.ReadWriterCloser that implements auto-reconnection.
func New(
	connect func(context.Context) (io.ReadWriteCloser, error),
) io.ReadWriteCloser {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &autoreconnector{
		connect:   connect,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
}

func (a *autoreconnector) Close() error {
	a.ctxCancel()

	a.connMutex.Lock()
	defer a.connMutex.Unlock()

	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	return nil
}

func (a *autoreconnector) getConnection(reset bool) (io.ReadWriteCloser, bool) {
	a.connMutex.Lock()
	defer a.connMutex.Unlock()

	if a.conn != nil {
		if !reset {
			return a.conn, true
		}

		a.conn.Close()
		a.conn = nil
	}

	select {
	case <-a.ctx.Done():
		return nil, false
	default:
	}

	for {
		var err error
		a.conn, err = a.connect(a.ctx)
		if err == nil {
			select {
			case <-a.ctx.Done():
				a.conn.Close()
				a.conn = nil
				return nil, false
			default:
			}

			return a.conn, true
		}

		select {
		case <-time.After(reconnectPeriod):
		case <-a.ctx.Done():
			return nil, false
		}
	}
}

func (a *autoreconnector) Read(p []byte) (int, error) {
	reset := false

	for {
		curConn, ok := a.getConnection(reset)
		if !ok {
			return 0, errTerminated
		}

		n, err := curConn.Read(p)
		if err == nil || (err == io.EOF && n > 0) {
			return n, err
		}

		reset = true
	}
}

func (a *autoreconnector) Write(p []byte) (int, error) {
	reset := false

	for {
		curConn, ok := a.getConnection(reset)
		if !ok {
			return 0, errTerminated
		}

		n, err := curConn.Write(p)
		if err == nil {
			return n, err
		}

		reset = true
	}
}
