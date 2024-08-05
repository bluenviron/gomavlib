package timednetconn

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeConn struct {
	readDeadlineDone  chan struct{}
	writeDeadlineDone chan struct{}
	closeDone         chan struct{}
}

func (fc *fakeConn) Read(_ []byte) (n int, err error) {
	return 10, nil
}

func (fc *fakeConn) Write(_ []byte) (n int, err error) {
	return 10, nil
}

func (fc *fakeConn) Close() error {
	close(fc.closeDone)
	return nil
}

func (fc *fakeConn) LocalAddr() net.Addr {
	return nil
}

func (fc *fakeConn) RemoteAddr() net.Addr {
	return nil
}

func (fc *fakeConn) SetDeadline(_ time.Time) error {
	return nil
}

func (fc *fakeConn) SetReadDeadline(_ time.Time) error {
	close(fc.readDeadlineDone)
	return nil
}

func (fc *fakeConn) SetWriteDeadline(_ time.Time) error {
	close(fc.writeDeadlineDone)
	return nil
}

func TestConn(t *testing.T) {
	fc := &fakeConn{
		readDeadlineDone:  make(chan struct{}),
		writeDeadlineDone: make(chan struct{}),
		closeDone:         make(chan struct{}),
	}

	conn := New(10*time.Second, 10*time.Second, fc)

	n, err := conn.Read(nil)
	require.NoError(t, err)
	require.Equal(t, 10, n)
	<-fc.readDeadlineDone

	n, err = conn.Write(nil)
	require.NoError(t, err)
	require.Equal(t, 10, n)
	<-fc.writeDeadlineDone

	err = conn.Close()
	require.NoError(t, err)
	<-fc.closeDone
}
