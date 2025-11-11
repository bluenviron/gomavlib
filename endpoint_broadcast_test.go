package gomavlib

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

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

	dialectRW := &dialect.ReadWriter{Dialect: testDialect}
	err = dialectRW.Initialize()
	require.NoError(t, err)

	rw := &frame.ReadWriter{
		ByteReadWriter: &readWriterFromFuncs{
			readFunc: func(p []byte) (int, error) {
				n, _, err2 := pc.ReadFrom(p)
				return n, err2
			},
			writeFunc: func(p []byte) (int, error) {
				return pc.WriteTo(p, &net.UDPAddr{
					IP:   net.ParseIP("127.255.255.255"),
					Port: 5601,
				})
			},
		},
		DialectRW: dialectRW,
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
