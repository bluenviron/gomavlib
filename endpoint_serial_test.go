package gomavlib

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

var _ endpointChannelSingle = (*endpointSerial)(nil)

func TestEndpointSerial(t *testing.T) {
	endpointCreated := make(chan *dummyEndpoint, 1)
	serialOpenFunc = func(name string, baud int) (io.ReadWriteCloser, error) {
		de := newDummyEndpoint()
		endpointCreated <- de
		return de, nil
	}

	dial := &dialect.Dialect{
		Version:  3,
		Messages: []message.Message{&MessageHeartbeat{}},
	}

	node, err := NewNode(NodeConf{
		Dialect:     dial,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{EndpointSerial{
			Device: "/dev/ttyUSB0",
			Baud:   57600,
		}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	<-endpointCreated

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	de := <-endpointCreated

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

func TestEndpointSerialReconnect(t *testing.T) {
	endpointCreated := make(chan *dummyEndpoint, 1)
	serialOpenFunc = func(name string, baud int) (io.ReadWriteCloser, error) {
		de := newDummyEndpoint()
		endpointCreated <- de
		return de, nil
	}

	dial := &dialect.Dialect{
		Version:  3,
		Messages: []message.Message{&MessageHeartbeat{}},
	}

	node, err := NewNode(NodeConf{
		Dialect:     dial,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{EndpointSerial{
			Device: "/dev/ttyUSB0",
			Baud:   57600,
		}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	<-endpointCreated

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	de := <-endpointCreated

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

	for i := 0; i < 2; i++ {
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

		evt := <-node.Events()
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
	}

	de.simulateReadError()
	de = <-endpointCreated

	for i := 0; i < 2; i++ {
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
		de.chOut <- buf.Bytes()
		buf.Reset()

		evt := <-node.Events()
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceID:  2 + byte(i),
				SystemID:    11,
				ComponentID: 1,
				Message:     msg,
				Checksum:    evt.(*EventFrame).Frame.GetChecksum(),
			},
			Channel: evt.(*EventFrame).Channel,
		}, evt)
	}
}
