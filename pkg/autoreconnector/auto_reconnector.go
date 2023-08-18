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
	errReconnecting = errors.New("reconnecting")
)

type state int

const (
	stateInitial state = iota
	stateReconnecting
	stateConnected
	stateTerminated
)

type autoReconnector struct {
	connect func(context.Context) (io.ReadWriteCloser, error)

	mutex            sync.Mutex
	state            state
	conn             io.ReadWriteCloser
	connectCtx       context.Context
	connectCtxCancel func()
}

// New returns a io.ReadWriterCloser that implements auto-reconnection.
func New(
	connect func(context.Context) (io.ReadWriteCloser, error),
) io.ReadWriteCloser {
	a := &autoReconnector{
		connect: connect,
	}

	a.resetConnection()

	return a
}

func (a *autoReconnector) Close() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.state = stateTerminated

	if a.connectCtxCancel != nil {
		a.connectCtxCancel()
	}

	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	return nil
}

func (a *autoReconnector) getConnection() (io.ReadWriteCloser, context.Context, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	switch a.state {
	case stateTerminated:
		return nil, nil, errTerminated

	case stateReconnecting:
		return nil, a.connectCtx, errReconnecting

	default:
		return a.conn, nil, nil
	}
}

func (a *autoReconnector) resetConnection() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	switch a.state {
	case stateTerminated, stateReconnecting:
		return
	}

	a.state = stateReconnecting

	if a.conn != nil {
		a.conn.Close()
		a.conn = nil
	}

	a.connectCtx, a.connectCtxCancel = context.WithCancel(context.Background())

	go func() {
		for {
			newConn, err := a.connect(a.connectCtx)
			if err == nil {
				a.setConn(newConn)
				return
			}

			select {
			case <-time.After(reconnectPeriod):
			case <-a.connectCtx.Done():
				return
			}
		}
	}()
}

func (a *autoReconnector) setConn(newConn io.ReadWriteCloser) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.state != stateReconnecting {
		newConn.Close()
		return
	}

	a.connectCtxCancel()
	a.connectCtxCancel = nil

	a.conn = newConn
	a.state = stateConnected
}

func (a *autoReconnector) Read(p []byte) (int, error) {
	for {
		curConn, connectCtx, err := a.getConnection()
		if err == errReconnecting {
			<-connectCtx.Done()
			continue
		}
		if err != nil {
			return 0, err
		}

		n, err := curConn.Read(p)

		if n == 0 {
			a.resetConnection()
			continue
		}

		return n, err
	}
}

func (a *autoReconnector) Write(p []byte) (int, error) {
	curConn, _, err := a.getConnection()
	if err != nil {
		return 0, err
	}

	n, err := curConn.Write(p)

	if n == 0 {
		a.resetConnection()
		return 0, errReconnecting
	}

	return n, err
}
