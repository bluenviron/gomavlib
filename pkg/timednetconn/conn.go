// Package timednetconn contains a net.Conn wrapper with deadlines.
package timednetconn

import (
	"io"
	"net"
	"time"
)

type conn struct {
	readTimeout  time.Duration
	writeTimeout time.Duration
	wrapped      net.Conn
}

// New returns a io.ReadWriteCloser that calls SetReadDeadline() before Read()
// and SetWriteDeadline() before Write().
func New(
	readTimeout time.Duration,
	writeTimeout time.Duration,
	wrapped net.Conn,
) io.ReadWriteCloser {
	return &conn{
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		wrapped:      wrapped,
	}
}

func (c *conn) Close() error {
	return c.wrapped.Close()
}

func (c *conn) Read(buf []byte) (int, error) {
	err := c.wrapped.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return 0, err
	}
	return c.wrapped.Read(buf)
}

func (c *conn) Write(buf []byte) (int, error) {
	err := c.wrapped.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if err != nil {
		return 0, err
	}
	return c.wrapped.Write(buf)
}
