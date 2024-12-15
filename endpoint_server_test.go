package gomavlib

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

var _ endpointChannelProvider = (*endpointServer)(nil)

func TestEndpointServer(t *testing.T) {
	for _, ca := range []string{"tcp", "udp"} {
		t.Run(ca, func(t *testing.T) {
			var e EndpointConf
			if ca == "tcp" {
				e = EndpointTCPServer{"127.0.0.1:5601"}
			} else {
				e = EndpointUDPServer{"127.0.0.1:5601"}
			}

			node, err := NewNode(NodeConf{
				Dialect:          testDialect,
				OutVersion:       V2,
				OutSystemID:      10,
				Endpoints:        []EndpointConf{e},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer node.Close()

			conn, err := net.Dial(ca, "127.0.0.1:5601")
			require.NoError(t, err)
			defer conn.Close()

			dialectRW, err := dialect.NewReadWriter(testDialect)
			require.NoError(t, err)

			rw, err := frame.NewReadWriter(frame.ReadWriterConf{
				ReadWriter:  conn,
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

				if i == 0 {
					evt := <-node.Events()
					require.Equal(t, &EventChannelOpen{
						Channel: evt.(*EventChannelOpen).Channel,
					}, evt)
				}

				evt := <-node.Events()
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
		})
	}
}

func TestEndpointServerIdleTimeout(t *testing.T) {
	for _, ca := range []string{"tcp", "udp"} {
		t.Run(ca, func(t *testing.T) {
			var e EndpointConf
			if ca == "tcp" {
				e = EndpointTCPServer{"127.0.0.1:5601"}
			} else {
				e = EndpointUDPServer{"127.0.0.1:5601"}
			}

			node, err := NewNode(NodeConf{
				Dialect:          testDialect,
				OutVersion:       V2,
				OutSystemID:      10,
				Endpoints:        []EndpointConf{e},
				HeartbeatDisable: true,
				IdleTimeout:      500 * time.Millisecond,
			})
			require.NoError(t, err)
			defer node.Close()

			conn, err := net.Dial(ca, "127.0.0.1:5601")
			require.NoError(t, err)
			defer conn.Close()

			dialectRW, err := dialect.NewReadWriter(testDialect)
			require.NoError(t, err)

			rw, err := frame.NewReadWriter(frame.ReadWriterConf{
				ReadWriter:  conn,
				DialectRW:   dialectRW,
				OutVersion:  frame.V2,
				OutSystemID: 11,
			})
			require.NoError(t, err)

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

			evt := <-node.Events()
			ch := evt.(*EventChannelOpen).Channel
			require.Equal(t, &EventChannelOpen{
				Channel: ch,
			}, evt)

			// frame
			<-node.Events()

			evt = <-node.Events()
			require.Equal(t, &EventChannelClose{
				Channel: ch,
			}, evt)
		})
	}
}
