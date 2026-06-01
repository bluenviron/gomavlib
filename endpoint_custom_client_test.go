package gomavlib

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
	"github.com/bluenviron/gomavlib/v4/pkg/frame"
	"github.com/bluenviron/gomavlib/v4/pkg/streamwriter"
	"github.com/stretchr/testify/require"
)

type dummyReadWriter struct {
	chOut     chan []byte
	chIn      chan []byte
	chReadErr chan struct{}
}

func newDummyReadWriterPair() (*dummyReadWriter, *dummyReadWriter) {
	one := &dummyReadWriter{
		chOut:     make(chan []byte),
		chIn:      make(chan []byte),
		chReadErr: make(chan struct{}),
	}

	two := &dummyReadWriter{
		chOut:     one.chIn,
		chIn:      one.chOut,
		chReadErr: make(chan struct{}),
	}

	return one, two
}

func (e *dummyReadWriter) simulateReadError() {
	close(e.chReadErr)
}

func (e *dummyReadWriter) Close() error {
	close(e.chOut)
	close(e.chIn)
	return nil
}

func (e *dummyReadWriter) Read(p []byte) (int, error) {
	select {
	case buf, ok := <-e.chOut:
		if !ok {
			return 0, io.EOF
		}
		return copy(p, buf), nil
	case <-e.chReadErr:
		return 0, errors.New("custom error")
	}
}

func (e *dummyReadWriter) Write(p []byte) (int, error) {
	e.chIn <- p
	return len(p), nil
}

func TestEndpointCustomClient(t *testing.T) {
	remote, local := newDummyReadWriterPair()

	node := &Node{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		HeartbeatDisable: true,
		Endpoints: []Endpoint{&EndpointCustomClient{
			Connect: func(_ context.Context) (net.Conn, error) {
				return &rwcToConn{remote}, nil
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
		ByteReadWriter: local,
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

type dummyReadWriter2 struct {
	terminate chan struct{}
}

func (r *dummyReadWriter2) Close() error {
	close(r.terminate)
	return nil
}

func (r *dummyReadWriter2) Read(_ []byte) (int, error) {
	<-r.terminate
	return 0, io.EOF
}

var errDummy = fmt.Errorf("dummy error")

func (r *dummyReadWriter2) Write(_ []byte) (int, error) {
	return 0, errDummy
}

func TestEndpointCustomClientWriteError(t *testing.T) {
	rwc := &dummyReadWriter2{
		terminate: make(chan struct{}),
	}

	node := &Node{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		HeartbeatDisable: true,
		Endpoints: []Endpoint{&EndpointCustomClient{
			Connect: func(_ context.Context) (net.Conn, error) {
				return &rwcToConn{rwc}, nil
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

	err = node.WriteMessageAll(&MessageHeartbeat{
		Type:           1,
		Autopilot:      2,
		BaseMode:       3,
		CustomMode:     6,
		SystemStatus:   4,
		MavlinkVersion: 5,
	})
	require.NoError(t, err)

	evt = <-node.Events()
	require.Equal(t, &EventChannelClose{
		Channel: evt.(*EventChannelClose).Channel,
		Error:   errDummy,
	}, evt)
}
