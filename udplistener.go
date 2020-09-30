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

type udpListenerConnIndex struct {
	IP   [4]byte
	Port int
}

type udpListenerConn struct {
	listener      *udpListener
	index         udpListenerConnIndex
	addr          *net.UDPAddr
	closed        bool
	readDeadline  time.Time
	writeDeadline time.Time

	read chan []byte
}

func newUdpListenerConn(listener *udpListener, index udpListenerConnIndex, addr *net.UDPAddr) *udpListenerConn {
	return &udpListenerConn{
		listener: listener,
		index:    index,
		addr:     addr,
		read:     make(chan []byte),
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
	c.listener.readMutex.Lock()
	defer c.listener.readMutex.Unlock()

	if c.closed == true {
		return nil
	}

	c.closed = true
	delete(c.listener.conns, c.index)

	// release anyone waiting on Read()
	close(c.read)

	// close socket when both listener and connections are closed
	if c.listener.closed == true && len(c.listener.conns) == 0 {
		c.listener.packetConn.Close()
	}

	return nil
}

// read synchronously, such that buffer can be freed after reading
func (c *udpListenerConn) Read(byt []byte) (int, error) {
	var buf []byte
	var ok bool

	if !c.readDeadline.IsZero() {
		readTimer := time.NewTimer(c.readDeadline.Sub(time.Now()))
		defer readTimer.Stop()

		select {
		case <-readTimer.C:
			return 0, udpErrorTimeout
		case buf, ok = <-c.read:
		}
	} else {
		buf, ok = <-c.read
	}

	if !ok {
		return 0, udpErrorTerminated
	}

	copy(byt, buf)
	c.listener.readDone <- struct{}{}
	return len(buf), nil
}

// write synchronously, such that buffer can be freed after writing
func (c *udpListenerConn) Write(byt []byte) (int, error) {
	c.listener.writeMutex.Lock()
	defer c.listener.writeMutex.Unlock()

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
	conns      map[udpListenerConnIndex]*udpListenerConn
	readMutex  sync.Mutex
	writeMutex sync.Mutex
	closed     bool

	acceptc  chan net.Conn
	readDone chan struct{}
}

func newUdpListener(network, address string) (net.Listener, error) {
	packetConn, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	l := &udpListener{
		packetConn: packetConn,
		conns:      make(map[udpListenerConnIndex]*udpListenerConn),
		acceptc:    make(chan net.Conn),
		readDone:   make(chan struct{}),
	}

	go l.reader()

	return l, nil
}

func (l *udpListener) Close() error {
	l.readMutex.Lock()
	defer l.readMutex.Unlock()

	if l.closed == true {
		return nil
	}

	l.closed = true

	// release anyone waiting on Accept()
	close(l.acceptc)

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

		// use ip and port as connection index
		uaddr := addr.(*net.UDPAddr)
		connIndex := udpListenerConnIndex{}
		connIndex.Port = uaddr.Port
		copy(connIndex.IP[:], uaddr.IP)

		func() {
			l.readMutex.Lock()
			defer l.readMutex.Unlock()

			conn, preExisting := l.conns[connIndex]

			if !preExisting && l.closed == true {
				// listener is closed, ignore new connection

			} else {
				if !preExisting {
					conn = newUdpListenerConn(l, connIndex, uaddr)
					l.conns[connIndex] = conn
					l.acceptc <- conn
				}

				// route buffer to connection
				conn.read <- buf[:n]

				// wait copy since buffer is shared
				<-l.readDone
			}
		}()
	}
}

func (l *udpListener) Accept() (net.Conn, error) {
	conn, ok := <-l.acceptc
	if !ok {
		return nil, udpErrorTerminated
	}
	return conn, nil
}
