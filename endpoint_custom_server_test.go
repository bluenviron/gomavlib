package gomavlib

import (
	"errors"
	"net"
	"sync"
	"testing"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
	"github.com/bluenviron/gomavlib/v4/pkg/frame"
	"github.com/bluenviron/gomavlib/v4/pkg/streamwriter"
	"github.com/stretchr/testify/require"
)

type mockNetListener struct {
	connCh    chan net.Conn
	terminate chan struct{}
	closeOnce sync.Once
}

func (l *mockNetListener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.terminate:
		return nil, errors.New("listener closed")
	}
}

func (l *mockNetListener) Close() error {
	l.closeOnce.Do(func() {
		close(l.terminate)
	})
	return nil
}

func (l *mockNetListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

func TestEndpointCustomServer(t *testing.T) {
	serverConn, clientConn := net.Pipe()

	ln := &mockNetListener{
		connCh:    make(chan net.Conn, 1),
		terminate: make(chan struct{}),
	}
	ln.connCh <- serverConn

	node := &Node{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		HeartbeatDisable: true,
		Endpoints: []Endpoint{&EndpointCustomServer{
			Listen: func() (net.Listener, error) {
				return ln, nil
			},
		}},
	}
	err := node.Initialize()
	require.NoError(t, err)
	defer node.Close()

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	dialectRW := &dialect.ReadWriter{Dialect: testDialect}
	err = dialectRW.Initialize()
	require.NoError(t, err)

	rw := &frame.ReadWriter{
		ByteReadWriter: clientConn,
		DialectRW:      dialectRW,
	}
	err = rw.Initialize()
	require.NoError(t, err)

	sw := &streamwriter.Writer{
		FrameWriter: rw.Writer,
		Version:     streamwriter.V2,
		SystemID:    11,
	}
	err = sw.Initialize()
	require.NoError(t, err)

	for i := range 3 { //nolint:dupl
		msg := &MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     6,
			SystemStatus:   4,
			MavlinkVersion: 5,
		}
		err = sw.Write(msg)
		require.NoError(t, err)

		evt = <-node.Events()
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceNumber: byte(i),
				SystemID:       11,
				ComponentID:    1,
				Message:        msg,
				Checksum:       evt.(*EventFrame).Frame.GetChecksum(),
			},
			Channel: evt.(*EventFrame).Channel,
		}, evt)

		msg = &MessageHeartbeat{
			Type:           6,
			Autopilot:      5,
			BaseMode:       4,
			CustomMode:     3,
			SystemStatus:   2,
			MavlinkVersion: 1,
		}
		err = node.WriteMessageAll(msg)
		require.NoError(t, err)

		var fr frame.Frame
		fr, err = rw.Read()
		require.NoError(t, err)
		require.Equal(t, &frame.V2Frame{
			SequenceNumber: byte(i),
			SystemID:       10,
			ComponentID:    1,
			Message:        msg,
			Checksum:       fr.GetChecksum(),
		}, fr)
	}
}
