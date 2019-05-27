package gomavlib

import (
	"net"
	"sync"
	"time"
)

// this file provides a net.Listener for udp servers, such that they can be
// handled like tcp ones.

// implements net.Error
type udpNetError struct {
	str       string
	isTimeout bool
}

func (e udpNetError) Error() string {
	return e.str
}

func (e udpNetError) Timeout() bool {
	return e.isTimeout
}

func (udpNetError) Temporary() bool {
	return false
}

var udpErrorTimeout net.Error = udpNetError{"timeout", true}
var udpErrorTerminated net.Error = udpNetError{"terminated", false}

type udpListenerConn struct {
	listener      *udpListener
	addr          net.Addr
	readChan      chan []byte
	closed        bool
	readDeadline  time.Time
	writeDeadline time.Time
}

func newUdpListenerConn(listener *udpListener, addr net.Addr) *udpListenerConn {
	return &udpListenerConn{
		listener: listener,
		addr:     addr,
		readChan: make(chan []byte),
	}
}

func (c *udpListenerConn) LocalAddr() net.Addr {
	// not implemented
	return nil
}

func (c *udpListenerConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *udpListenerConn) Close() error {
	c.listener.mutex.Lock()
	defer c.listener.mutex.Unlock()

	if c.closed == true {
		return nil
	}

	c.closed = true
	delete(c.listener.conns, c.addr.String())

	// release anyone waiting on Read()
	close(c.readChan)

	// close socket when both listener and connections are closed
	if c.listener.closed == true && len(c.listener.conns) == 0 {
		c.listener.packetConn.Close()
	}

	return nil
}

func (c *udpListenerConn) Read(byt []byte) (int, error) {
	var buf []byte
	var ok bool

	if !c.readDeadline.IsZero() {
		readTimer := time.NewTimer(c.readDeadline.Sub(time.Now()))
		defer readTimer.Stop()

		select {
		case <-readTimer.C:
			return 0, udpErrorTimeout
		case buf, ok = <-c.readChan:
		}
	} else {
		buf, ok = <-c.readChan
	}

	if ok == false {
		return 0, udpErrorTerminated
	}

	copy(byt, buf)
	c.listener.readDone <- struct{}{}
	return len(buf), nil
}

// write synchronously, such that buffer can be freed after writing
func (c *udpListenerConn) Write(byt []byte) (int, error) {
	c.listener.mutex.Lock()
	defer c.listener.mutex.Unlock()

	if c.closed == true {
		return 0, udpErrorTerminated
	}

	if !c.writeDeadline.IsZero() {
		err := c.listener.packetConn.SetWriteDeadline(c.writeDeadline)
		if err != nil {
			return 0, err
		}
	}

	return c.listener.packetConn.WriteTo(byt, c.addr)
}

func (c *udpListenerConn) SetDeadline(time.Time) error {
	// not implemented
	return nil
}

func (c *udpListenerConn) SetReadDeadline(t time.Time) error {
	c.readDeadline = t
	return nil
}

func (c *udpListenerConn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}

type udpListener struct {
	packetConn net.PacketConn
	acceptChan chan net.Conn
	readDone   chan struct{}
	conns      map[string]*udpListenerConn
	mutex      sync.Mutex
	closed     bool
}

func newUdpListener(network, address string) (net.Listener, error) {
	packetConn, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	l := &udpListener{
		packetConn: packetConn,
		acceptChan: make(chan net.Conn),
		readDone:   make(chan struct{}),
		conns:      make(map[string]*udpListenerConn),
	}

	go l.reader()

	return l, nil
}

func (l *udpListener) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed == true {
		return nil
	}

	l.closed = true

	// release anyone waiting on Accept()
	close(l.acceptChan)

	// close socket when both listener and connections are closed
	if len(l.conns) == 0 {
		l.packetConn.Close()
	}

	return nil
}

func (l *udpListener) Addr() net.Addr {
	return l.packetConn.LocalAddr()
}

func (l *udpListener) reader() {
	buf := make([]byte, 2048) // MTU is ~1500

	for {
		// read WITHOUT deadline. Long periods without packets are normal since
		// we're not directly connected to someone.
		n, addr, err := l.packetConn.ReadFrom(buf)
		if err != nil {
			break
		}

		// addr changes every time ReadFrom() is called, we use its string
		// representation as index
		addrs := addr.String()

		func() {
			l.mutex.Lock()
			defer l.mutex.Unlock()

			conn, preExisting := l.conns[addrs]

			if preExisting == false && l.closed == true {
				// listener is closed, ignore new connection

			} else {
				if preExisting == false {
					conn = newUdpListenerConn(l, addr)
					l.conns[addrs] = conn
					l.acceptChan <- conn
				}

				// route buffer to connection
				conn.readChan <- buf[:n]

				// wait copy since buffer is shared
				<-l.readDone
			}
		}()
	}
}

func (l *udpListener) Accept() (net.Conn, error) {
	conn, ok := <-l.acceptChan
	if ok == false {
		return nil, udpErrorTerminated
	}
	return conn, nil
}
