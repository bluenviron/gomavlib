// Package timednetconn contains a net.Conn wrapper that calls SetWriteDeadline() before Write().
package timednetconn

import (
	"io"
	"net"
	"time"
)

type timednetconn struct {
	writeTimeout time.Duration
	wrapped      net.Conn
}

// New returns a io.ReadWriteCloser that calls SetWriteDeadline() before Write().
func New(writeTimeout time.Duration, wrapped net.Conn) io.ReadWriteCloser {
	return &timednetconn{
		writeTimeout: writeTimeout,
		wrapped:      wrapped,
	}
}

func (c *timednetconn) Close() error {
	return c.wrapped.Close()
}

func (c *timednetconn) Read(buf []byte) (int, error) {
	// do not call SetReadDeadline()
	// since we don't want to disconnect in case of long pauses between messages
	return c.wrapped.Read(buf)
}

func (c *timednetconn) Write(buf []byte) (int, error) {
	err := c.wrapped.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return 0, err
	}
	return c.wrapped.Write(buf)
}
