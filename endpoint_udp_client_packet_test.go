package gomavlib

import (
	"bytes"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

// buildValidDatagram builds a serialised MAVLink v2 frame for msg via the
// test dialect and returns the raw bytes suitable for sending as one UDP datagram.
func buildValidDatagram(t *testing.T, msg *MessageHeartbeat, systemID byte) []byte {
	t.Helper()

	var buf bytes.Buffer
	drw := &dialect.ReadWriter{Dialect: testDialect}
	require.NoError(t, drw.Initialize())

	fw := &frame.Writer{ByteWriter: &buf, DialectRW: drw}
	require.NoError(t, fw.Initialize())

	sw := &streamwriter.Writer{
		FrameWriter: fw,
		Version:     streamwriter.V2,
		SystemID:    systemID,
	}
	require.NoError(t, sw.Initialize())
	require.NoError(t, sw.Write(msg))

	return append([]byte(nil), buf.Bytes()...)
}

// TestEndpointUDPClientPacketRecoverAfterMalformedDatagrams verifies that
// EndpointUDPClientPacket discards each bad datagram atomically and resumes
// normal parsing on the next one, without cascading spurious parse errors.
func TestEndpointUDPClientPacketRecoverAfterMalformedDatagrams(t *testing.T) {
	pc, err := net.ListenPacket("udp4", "127.0.0.1:5605")
	require.NoError(t, err)
	defer pc.Close()

	validData := buildValidDatagram(t, &MessageHeartbeat{
		Type:           6,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}, 11)

	// Two distinct malformed datagrams.
	junk1 := []byte("\x11\x22\x33\x44\x55\x66\x77\x88")
	junk2 := []byte("\xAA\xBB\xCC\xDD\xEE\xFF\x01\x02")

	serverErr := make(chan error, 1)

	go func() {
		buf := make([]byte, 2048)
		n, addr, err2 := pc.ReadFrom(buf)
		if err2 != nil {
			serverErr <- err2
			return
		}
		if n <= 0 {
			serverErr <- fmt.Errorf("empty initial datagram")
			return
		}

		for _, junk := range [][]byte{junk1, junk2} {
			if _, err2 = pc.WriteTo(junk, addr); err2 != nil {
				serverErr <- err2
				return
			}
		}

		if _, err2 = pc.WriteTo(validData, addr); err2 != nil {
			serverErr <- err2
			return
		}

		serverErr <- nil
	}()

	node, err := NewNode(NodeConf{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointUDPClientPacket{"127.0.0.1:5605"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	// Wait for channel open.
	evt := <-node.Events()
	_, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Send one frame to trigger server goroutine.
	err = node.WriteMessageAll(&MessageHeartbeat{
		Type: 1, Autopilot: 2, BaseMode: 3,
		CustomMode: 6, SystemStatus: 4, MavlinkVersion: 5,
	})
	require.NoError(t, err)

	deadline := time.After(2 * time.Second)
	parseErrors := 0

	for {
		select {
		case sErr := <-serverErr:
			require.NoError(t, sErr)

		case ev := <-node.Events():
			switch et := ev.(type) {
			case *EventParseError:
				parseErrors++

			case *EventFrame:
				// Must have received exactly 2 parse errors (one per junk datagram).
				require.Equal(t, 2, parseErrors,
					"expected exactly 2 EventParseError (one per junk datagram), got %d", parseErrors)

				require.Equal(t, &frame.V2Frame{
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
					Checksum: et.Frame.GetChecksum(),
				}, et.Frame)
				return
			}

		case <-deadline:
			t.Fatalf("timed out: got %d parse errors, never received valid frame", parseErrors)
		}
	}
}
