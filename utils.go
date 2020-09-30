package gomavlib

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

var errorTerminated = fmt.Errorf("terminated")

// netTimedConn forces a net.Conn to use timeouts
type netTimedConn struct {
	conn net.Conn
}

func (c *netTimedConn) Close() error {
	return c.conn.Close()
}

func (c *netTimedConn) Read(buf []byte) (int, error) {
	err := c.conn.SetReadDeadline(time.Now().Add(netReadTimeout))
	if err != nil {
		return 0, err
	}
	return c.conn.Read(buf)
}

func (c *netTimedConn) Write(buf []byte) (int, error) {
	err := c.conn.SetWriteDeadline(time.Now().Add(netWriteTimeout))
	if err != nil {
		return 0, err
	}
	return c.conn.Write(buf)
}

func randomByte() byte {
	var buf [1]byte
	rand.Read(buf[:])
	return buf[0]
}
