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

var errorTimeout net.Error = udpNetError{"timeout", true}
var errorTerminated net.Error = udpNetError{"terminated", false}

type udpListenerConn struct {
	listener      *udpListener
	addr          net.Addr
	readChan      chan []byte
	closed        bool
	readDeadline  time.Time
	writeDeadline time.Time
}

func (c *udpListenerConn) Close() error {
	exit := false
	func() {
		c.listener.mutex.Lock()
		defer c.listener.mutex.Unlock()

		if c.closed == true {
			exit = true
			return
		}
		c.closed = true

		delete(c.listener.conns, c.addr.String())
	}()
	if exit == true {
		return nil
	}

	// release anyone waiting on Read()
	close(c.readChan)

	return nil
}

func (c *udpListenerConn) LocalAddr() net.Addr {
	// not implemented
	return nil
}

func (c *udpListenerConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *udpListenerConn) Read(byt []byte) (int, error) {
	exit := false
	func() {
		c.listener.mutex.Lock()
		defer c.listener.mutex.Unlock()

		if c.closed == true {
			exit = true
			return
		}
	}()
	if exit == true {
		return 0, errorTerminated
	}

	readTimer := time.NewTimer(c.readDeadline.Sub(time.Now()))
	defer readTimer.Stop()

	select {
	case <-readTimer.C:
		return 0, errorTimeout

	case buf, ok := <-c.readChan:
		if ok == false {
			return 0, errorTerminated
		}
		copy(byt, buf)
		return len(buf), nil
	}
}

func (c *udpListenerConn) Write(byt []byte) (int, error) {
	c.listener.mutex.Lock()
	defer c.listener.mutex.Unlock()

	if c.closed == true {
		return 0, errorTerminated
	}

	if c.listener.closed == true {
		return 0, errorTerminated
	}

	// write synchronously, such that buffer can be freed after writing
	c.listener.pc.SetWriteDeadline(c.writeDeadline)
	return c.listener.pc.WriteTo(byt, c.addr)
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
	pc         net.PacketConn
	done       chan struct{}
	acceptChan chan net.Conn
	conns      map[string]*udpListenerConn
	mutex      sync.Mutex
	closed     bool
}

func newUdpListener(network, address string) (net.Listener, error) {
	pc, err := net.ListenPacket(network, address)
	if err != nil {
		return nil, err
	}

	l := &udpListener{
		pc:         pc,
		done:       make(chan struct{}),
		acceptChan: make(chan net.Conn),
		conns:      make(map[string]*udpListenerConn),
	}
	go l.do()
	return l, nil
}

func (l *udpListener) Close() error {
	exit := false
	func() {
		l.mutex.Lock()
		defer l.mutex.Unlock()

		if l.closed == true {
			exit = true
			return
		}
		l.closed = true
	}()
	if exit == true {
		return nil
	}

	// terminate do()
	l.pc.Close()

	// wait do() and ensure that acceptChan is not used anymore
	<-l.done

	// release anyone waiting on Accept()
	close(l.acceptChan)

	return nil
}

func (l *udpListener) Addr() net.Addr {
	return l.pc.LocalAddr()
}

func (l *udpListener) do() {
	defer func() { l.done <- struct{}{} }()

	for {
		// buffer is be passed to the routines so it cannot be shared
		buf := make([]byte, 2048) // MTU is 1500

		n, addr, err := l.pc.ReadFrom(buf)
		if err != nil {
			break
		}
		buf = buf[:n]

		// addr changes every time ReadFrom() is called, we use its string
		// representation as index
		addrs := addr.String()

		var conn *udpListenerConn
		var preExisting bool

		func() {
			l.mutex.Lock()
			defer l.mutex.Unlock()

			conn, preExisting = l.conns[addrs]
			if preExisting == false {
				conn = &udpListenerConn{
					listener: l,
					addr:     addr,
					readChan: make(chan []byte),
				}
				l.conns[addrs] = conn
			}
		}()

		if preExisting == false {
			l.acceptChan <- conn
		}

		// route buffer to connection
		conn.readChan <- buf
	}
}

func (l *udpListener) Accept() (net.Conn, error) {
	conn, ok := <-l.acceptChan
	if ok == false {
		return nil, errorTerminated
	}
	return conn, nil
}
