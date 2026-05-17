package gomavlib

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

func TestEndpointUDPClientDatagramRecovery(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:5604")
	require.NoError(t, err)
	defer pc.Close()

	serverDone := make(chan struct{})

	go func() {
		defer close(serverDone)

		buf := make([]byte, 4096)
		_, clientAddr, err2 := pc.ReadFrom(buf)
		require.NoError(t, err2)

		// first malformed packet (too short)
		_, err2 = pc.WriteTo([]byte{frame.V2MagicByte}, clientAddr)
		require.NoError(t, err2)

		// second malformed packet (unknown incompatibility flag, with trailing payload+checksum bytes
		_, err2 = pc.WriteTo([]byte{frame.V2MagicByte, 5, 0x04, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, clientAddr)
		require.NoError(t, err2)

		// valid packet, with extra bytes
		_, err2 = pc.WriteTo([]byte{
			0xfd, 0x09, 0x00, 0x00, 0x00, 0x0b, 0x01, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x06, 0x05,
			0x04, 0x02, 0x01, 0xb4, 0xde,
			0xff, 0xff, 0xff, 0xff, 0xff,
		}, clientAddr)
		require.NoError(t, err2)

		// valid packet
		_, err2 = pc.WriteTo([]byte{
			0xfd, 0x09, 0x00, 0x00, 0x00, 0x0b, 0x01, 0x00,
			0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x06, 0x05,
			0x04, 0x02, 0x01, 0xb4, 0xde,
		}, clientAddr)
		require.NoError(t, err2)
	}()

	node := &Node{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5604"}},
		HeartbeatDisable: true,
	}
	err = node.Initialize()
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
	parseErr, ok := evt.(*EventParseError)
	require.True(t, ok)
	require.EqualError(t, parseErr.Error, "packet is too short")

	evt = <-node.Events()
	parseErr, ok = evt.(*EventParseError)
	require.True(t, ok)
	require.EqualError(t, parseErr.Error, "unknown incompatibility flag: 4")

	evt = <-node.Events()
	parseErr, ok = evt.(*EventParseError)
	require.True(t, ok)
	require.EqualError(t, parseErr.Error, "skipped 7 bytes")

	evt = <-node.Events()
	fr, ok := evt.(*EventFrame)
	require.True(t, ok)
	require.Equal(t, &EventFrame{
		Frame: &frame.V2Frame{
			SequenceNumber: 0,
			SystemID:       11,
			ComponentID:    1,
			Message: &MessageHeartbeat{
				Type:           6,
				Autopilot:      5,
				BaseMode:       4,
				CustomMode:     3,
				SystemStatus:   2,
				MavlinkVersion: 1,
			},
			Checksum: fr.Frame.GetChecksum(),
		},
		Channel: fr.Channel,
	}, evt)

	evt = <-node.Events()
	parseErr, ok = evt.(*EventParseError)
	require.True(t, ok)
	require.EqualError(t, parseErr.Error, "skipped 5 bytes")

	evt = <-node.Events()
	fr, ok = evt.(*EventFrame)
	require.True(t, ok)
	require.Equal(t, &EventFrame{
		Frame: &frame.V2Frame{
			SequenceNumber: 0,
			SystemID:       11,
			ComponentID:    1,
			Message: &MessageHeartbeat{
				Type:           6,
				Autopilot:      5,
				BaseMode:       4,
				CustomMode:     3,
				SystemStatus:   2,
				MavlinkVersion: 1,
			},
			Checksum: fr.Frame.GetChecksum(),
		},
		Channel: fr.Channel,
	}, evt)

	<-serverDone
}
