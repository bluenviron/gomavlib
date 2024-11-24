package gomavlib

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/chrisdalke/gomavlib/v3/pkg/dialect"
	"github.com/chrisdalke/gomavlib/v3/pkg/frame"
)

var _ endpointChannelSingle = (*endpointCustom)(nil)

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

func TestEndpointCustom(t *testing.T) {
	remote, local := newDummyReadWriterPair()

	node, err := NewNode(NodeConf{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointCustom{remote}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	dialectRW, err := dialect.NewReadWriter(testDialect)
	require.NoError(t, err)

	rw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter:  local,
		DialectRW:   dialectRW,
		OutVersion:  frame.V2,
		OutSystemID: 11,
	})
	require.NoError(t, err)

	for i := 0; i < 3; i++ { //nolint:dupl
		msg := &MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     6,
			SystemStatus:   4,
			MavlinkVersion: 5,
		}
		err = rw.WriteMessage(msg)
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
		err := node.WriteMessageAll(msg)
		require.NoError(t, err)

		fr, err := rw.Read()
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
