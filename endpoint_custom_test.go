package gomavlib

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

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

	for i := 0; i < 3; i++ { //nolint:dupl
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
