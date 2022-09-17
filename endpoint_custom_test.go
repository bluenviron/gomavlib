package gomavlib

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/message"
)

var _ endpointChannelSingle = (*endpointCustom)(nil)

type customEndpoint struct {
	chOut chan []byte
	chIn  chan []byte
}

func (e *customEndpoint) Close() error {
	close(e.chOut)
	close(e.chIn)
	return nil
}

func (e *customEndpoint) Read(p []byte) (int, error) {
	buf, ok := <-e.chOut
	if !ok {
		return 0, io.EOF
	}
	return copy(p, buf), nil
}

func (e *customEndpoint) Write(p []byte) (int, error) {
	e.chIn <- p
	return len(p), nil
}

func TestEndpointCustom(t *testing.T) {
	dial := &dialect.Dialect{
		Version:  3,
		Messages: []message.Message{&MessageHeartbeat{}},
	}

	ce := &customEndpoint{
		chOut: make(chan []byte),
		chIn:  make(chan []byte),
	}

	node, err := NewNode(NodeConf{
		Dialect:          dial,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointCustom{ce}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	dialectRW, err := dialect.NewReadWriter(dial)
	require.NoError(t, err)

	var buf bytes.Buffer

	rw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter:  &buf,
		DialectRW:   dialectRW,
		OutVersion:  frame.V2,
		OutSystemID: 11,
	})
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
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

		ce.chOut <- buf.Bytes()
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

		buf2 := <-ce.chIn
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
