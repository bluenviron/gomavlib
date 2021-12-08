// Package udplistener provides a UDP-based Listener.
package udplistener

import (
	"net"
	"sync"
	"time"
)

// MTU is ~1500
const bufferSize = 2048

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

var (
	errTimeout    net.Error = udpNetError{"timeout", true}
	errTerminated net.Error = udpNetError{"terminated", false}
)

type connIndex struct {
	IP   [4]byte
	Port int
}

type conn struct {
	listener      *Listener
	index         connIndex
	addr          *net.UDPAddr
	closed        bool
	readDeadline  time.Time
	writeDeadline time.Time

	// in
	read chan []byte
}

func newConn(listener *Listener, index connIndex, addr *net.UDPAddr) *conn {
	return &conn{
		listener: listener,
		index:    index,
		addr:     addr,
		read:     make(chan []byte),
	}
}

// LocalAddr implements the net.Conn interface.
func (c *conn) LocalAddr() net.Addr {
	// not implemented
	return nil
}

// RemoteAddr implements the net.Conn interface.
func (c *conn) RemoteAddr() net.Addr {
	return c.addr
}

// Close implements the net.Conn interface.
func (c *conn) Close() error {
	c.listener.readMutex.Lock()
	defer c.listener.readMutex.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	delete(c.listener.conns, c.index)

	// release anyone waiting on Read()
	close(c.read)

	// close socket when both listener and connections are closed
	if c.listener.closed && len(c.listener.conns) == 0 {
		c.listener.pc.Close()
	}

	return nil
}

// Read implements the net.Conn interface.
func (c *conn) Read(byt []byte) (int, error) {
	var buf []byte
	var ok bool

	if !c.readDeadline.IsZero() {
		readTimer := time.NewTimer(time.Until(c.readDeadline))
		defer readTimer.Stop()

		select {
		case <-readTimer.C:
			return 0, errTimeout
		case buf, ok = <-c.read:
		}
	} else {
		buf, ok = <-c.read
	}

	if !ok {
		return 0, errTerminated
	}

	n := copy(byt, buf)
	c.listener.readDone <- n
	return n, nil
}

// Write implements the net.Conn interface.
func (c *conn) Write(byt []byte) (int, error) {
	c.listener.writeMutex.Lock()
	defer c.listener.writeMutex.Unlock()

	if !c.writeDeadline.IsZero() {
		err := c.listener.pc.SetWriteDeadline(c.writeDeadline)
		if err != nil {
			return 0, err
		}
	}

	return c.listener.pc.WriteTo(byt, c.addr)
}

// SetDeadline implements the net.Conn interface.
func (c *conn) SetDeadline(time.Time) error {
	// not implemented
	return nil
}

// SetReadDeadline implements the net.Conn interface.
func (c *conn) SetReadDeadline(t time.Time) error {
	c.readDeadline = t
	return nil
}

// SetWriteDeadline implements the net.Conn interface.
func (c *conn) SetWriteDeadline(t time.Time) error {
	c.writeDeadline = t
	return nil
}

// Listener is a UDP listener.
type Listener struct {
	pc         net.PacketConn
	conns      map[connIndex]*conn
	readMutex  sync.Mutex
	writeMutex sync.Mutex
	closed     bool

	accept   chan net.Conn
	readDone chan int
}

// New allocates a Listener.
func New(network, address string) (net.Listener, error) {
	pc, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	l := &Listener{
		pc:       pc,
		conns:    make(map[connIndex]*conn),
		accept:   make(chan net.Conn),
		readDone: make(chan int),
	}

	go l.runReader()

	return l, nil
}

// Close implements the net.Listener interface.
func (l *Listener) Close() error {
	l.readMutex.Lock()
	defer l.readMutex.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true

	// release anyone waiting on Accept()
	close(l.accept)

	// close socket when both listener and connections are closed
	if len(l.conns) == 0 {
		l.pc.Close()
	}

	return nil
}

// Addr implements the net.Listener interface.
func (l *Listener) Addr() net.Addr {
	return l.pc.LocalAddr()
}

func (l *Listener) runReader() {
	buf := make([]byte, bufferSize)

	for {
		// read WITHOUT deadline. Long periods without packets are normal since
		// we're not directly connected to someone.
		n, addr, err := l.pc.ReadFrom(buf)
		if err != nil {
			break
		}

		// use ip and port as connection index
		uaddr := addr.(*net.UDPAddr)
		connIndex := connIndex{}
		connIndex.Port = uaddr.Port
		copy(connIndex.IP[:], uaddr.IP)

		func() {
			l.readMutex.Lock()
			defer l.readMutex.Unlock()

			conn, preExisting := l.conns[connIndex]

			if !preExisting && l.closed {
				// listener is closed, ignore new connections
			} else {
				if !preExisting {
					conn = newConn(l, connIndex, uaddr)
					l.conns[connIndex] = conn
					l.accept <- conn
				}

				start := 0
				for n > 0 {
					// route buffer to connection
					conn.read <- buf[start : start+n]

					// wait copy since buffer is shared
					read := <-l.readDone
					n -= read
					start += read
				}
			}
		}()
	}
}

// Accept implements the net.Listener interface.
func (l *Listener) Accept() (net.Conn, error) {
	conn, ok := <-l.accept
	if !ok {
		return nil, errTerminated
	}
	return conn, nil
}
