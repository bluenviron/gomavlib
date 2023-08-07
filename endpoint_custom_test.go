package gomavlib

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
)

var _ endpointChannelSingle = (*endpointCustom)(nil)

type dummyEndpoint struct {
	chOut     chan []byte
	chIn      chan []byte
	chReadErr chan struct{}
}

func newDummyEndpoint() *dummyEndpoint {
	return &dummyEndpoint{
		chOut:     make(chan []byte),
		chIn:      make(chan []byte),
		chReadErr: make(chan struct{}),
	}
}

func (e *dummyEndpoint) simulateReadError() {
	close(e.chReadErr)
}

func (e *dummyEndpoint) push(buf []byte) {
	e.chOut <- buf
}

func (e *dummyEndpoint) pull() []byte {
	return <-e.chIn
}

func (e *dummyEndpoint) Close() error {
	close(e.chOut)
	close(e.chIn)
	return nil
}

func (e *dummyEndpoint) Read(p []byte) (int, error) {
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

func (e *dummyEndpoint) Write(p []byte) (int, error) {
	e.chIn <- p
	return len(p), nil
}

func TestEndpointCustom(t *testing.T) {
	de := newDummyEndpoint()

	node, err := NewNode(NodeConf{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointCustom{de}},
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

	var buf bytes.Buffer

	rw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter:  &buf,
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
		de.push(buf.Bytes())
		buf.Reset()

		evt = <-node.Events()
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceID:  byte(i),
				SystemID:    11,
				ComponentID: 1,
				Message:     msg,
				Checksum:    evt.(*EventFrame).Frame.GetChecksum(),
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
		node.WriteMessageAll(msg)

		buf2 := de.pull()
		buf.Write(buf2)
		fr, err := rw.Read()
		require.NoError(t, err)
		require.Equal(t, &frame.V2Frame{
			SequenceID:  byte(i),
			SystemID:    10,
			ComponentID: 1,
			Message:     msg,
			Checksum:    fr.GetChecksum(),
		}, fr)
	}
}
