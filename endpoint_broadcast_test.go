package gomavlib

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
)

var _ endpointChannelSingle = (*endpointUDPBroadcast)(nil)

type readWriterFromFuncs struct {
	readFunc  func([]byte) (int, error)
	writeFunc func([]byte) (int, error)
}

func (rw *readWriterFromFuncs) Read(p []byte) (int, error) {
	return rw.readFunc(p)
}

func (rw *readWriterFromFuncs) Write(p []byte) (int, error) {
	return rw.writeFunc(p)
}

func TestEndpointBroadcast(t *testing.T) {
	node, err := NewNode(NodeConf{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointUDPBroadcast{"127.255.255.255:5602", ":5601"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	pc, err := net.ListenPacket("udp4", ":5602")
	require.NoError(t, err)
	defer pc.Close()

	dialectRW, err := dialect.NewReadWriter(testDialect)
	require.NoError(t, err)

	rw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter: &readWriterFromFuncs{
			readFunc: func(p []byte) (int, error) {
				n, _, err := pc.ReadFrom(p)
				return n, err
			},
			writeFunc: func(p []byte) (int, error) {
				return pc.WriteTo(p, &net.UDPAddr{
					IP:   net.ParseIP("127.255.255.255"),
					Port: 5601,
				})
			},
		},
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
