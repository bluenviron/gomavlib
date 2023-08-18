// Package reconnector allows to perform automatic reconnections.
package reconnector

import (
	"context"
	"io"
	"sync"
	"time"
)

var reconnectPeriod = 2 * time.Second

// ConnectFunc is the prototype of the callback passed to New()
type ConnectFunc func(context.Context) (io.ReadWriteCloser, error)

type connWithContext struct {
	rwc       io.ReadWriteCloser
	mutex     sync.Mutex
	ctx       context.Context
	ctxCancel func()
}

func newConnWithContext(rwc io.ReadWriteCloser) *connWithContext {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &connWithContext{
		rwc:       rwc,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
}

func (c *connWithContext) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	select {
	case <-c.ctx.Done():
		return nil
	default:
	}

	c.ctxCancel()

	return c.rwc.Close()
}

func (c *connWithContext) Read(p []byte) (int, error) {
	n, err := c.rwc.Read(p)
	if n == 0 {
		c.Close() //nolint:errcheck
	}
	return n, err
}

func (c *connWithContext) Write(p []byte) (int, error) {
	n, err := c.rwc.Write(p)
	if n == 0 {
		c.Close() //nolint:errcheck
	}
	return n, err
}

// Reconnector allocws to perform automatic reconnections.
type Reconnector struct {
	connect ConnectFunc

	ctx       context.Context
	ctxCancel func()
	curConn   *connWithContext
}

// New allocates a Reconnector.
func New(connect ConnectFunc) *Reconnector {
	ctx, ctxCancel := context.WithCancel(context.Background())

	return &Reconnector{
		connect:   connect,
		ctx:       ctx,
		ctxCancel: ctxCancel,
	}
}

// Close closes a reconnector.
func (a *Reconnector) Close() {
	a.ctxCancel()
}

// Reconnect returns the next working connection.
func (a *Reconnector) Reconnect() (io.ReadWriteCloser, bool) {
	if a.curConn != nil {
		select {
		case <-a.curConn.ctx.Done():
		case <-a.ctx.Done():
			return nil, false
		}
	}

	for {
		conn, err := a.connect(a.ctx)
		if err != nil {
			select {
			case <-time.After(reconnectPeriod):
				continue
			case <-a.ctx.Done():
				return nil, false
			}
		}

		a.curConn = newConnWithContext(conn)
		return a.curConn, true
	}
}
