package gomavlib

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/pion/transport/v2/udp"
	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
)

var _ endpointChannelSingle = (*endpointClient)(nil)

func TestEndpointClient(t *testing.T) {
	for _, ca := range []string{"tcp", "udp"} {
		t.Run(ca, func(t *testing.T) {
			var ln net.Listener
			if ca == "tcp" {
				var err error
				ln, err = net.Listen("tcp", "127.0.0.1:5601")
				require.NoError(t, err)
			} else {
				addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:5601")
				require.NoError(t, err)

				ln, err = udp.Listen("udp", addr)
				require.NoError(t, err)
			}

			defer ln.Close()

			var e EndpointConf
			if ca == "tcp" {
				e = EndpointTCPClient{"127.0.0.1:5601"}
			} else {
				e = EndpointUDPClient{"127.0.0.1:5601"}
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

			evt := <-node.Events()
			require.Equal(t, &EventChannelOpen{
				Channel: evt.(*EventChannelOpen).Channel,
			}, evt)

			var rw *frame.ReadWriter

			for i := 0; i < 3; i++ {
				msg := &MessageHeartbeat{
					Type:           1,
					Autopilot:      2,
					BaseMode:       3,
					CustomMode:     6,
					SystemStatus:   4,
					MavlinkVersion: 5,
				}
				node.WriteMessageAll(msg)

				if i == 0 {
					conn, err := ln.Accept()
					require.NoError(t, err)
					defer conn.Close()

					dialectRW, err := dialect.NewReadWriter(testDialect)
					require.NoError(t, err)

					rw, err = frame.NewReadWriter(frame.ReadWriterConf{
						ReadWriter:  conn,
						DialectRW:   dialectRW,
						OutVersion:  frame.V2,
						OutSystemID: 11,
					})
					require.NoError(t, err)
				}

				fr, err := rw.Read()
				require.NoError(t, err)
				require.Equal(t, &frame.V2Frame{
					SequenceID:  byte(i),
					SystemID:    10,
					ComponentID: 1,
					Message:     msg,
					Checksum:    fr.GetChecksum(),
				}, fr)

				msg = &MessageHeartbeat{
					Type:           6,
					Autopilot:      5,
					BaseMode:       4,
					CustomMode:     3,
					SystemStatus:   2,
					MavlinkVersion: 1,
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
			}
		})
	}
}

func TestEndpointClientIdleTimeout(t *testing.T) {
	for _, ca := range []string{"tcp"} {
		t.Run(ca, func(t *testing.T) {
			var ln net.Listener
			var err error
			ln, err = net.Listen("tcp", "127.0.0.1:5603")
			require.NoError(t, err)

			defer ln.Close()

			var e EndpointConf
			if ca == "tcp" {
				e = EndpointTCPClient{"127.0.0.1:5603"}
			} else {
				e = EndpointUDPClient{"127.0.0.1:5603"}
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

			evt := <-node.Events()
			require.Equal(t, &EventChannelOpen{
				Channel: evt.(*EventChannelOpen).Channel,
			}, evt)

			msg := &MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     6,
				SystemStatus:   4,
				MavlinkVersion: 5,
			}
			node.WriteMessageAll(msg)

			conn, err := ln.Accept()
			require.NoError(t, err)

			dialectRW, err := dialect.NewReadWriter(testDialect)
			require.NoError(t, err)

			rw, err := frame.NewReadWriter(frame.ReadWriterConf{
				ReadWriter:  conn,
				DialectRW:   dialectRW,
				OutVersion:  frame.V2,
				OutSystemID: 11,
			})
			require.NoError(t, err)

			fr, err := rw.Read()
			require.NoError(t, err)
			require.Equal(t, &frame.V2Frame{
				SequenceID:  0,
				SystemID:    10,
				ComponentID: 1,
				Message:     msg,
				Checksum:    fr.GetChecksum(),
			}, fr)

			closed := make(chan struct{})

			go func() {
				_, err = rw.Read()
				require.Equal(t, io.EOF, err)
				conn.Close()
				close(closed)
			}()

			select {
			case <-closed:
			case <-time.After(1 * time.Second):
				t.Errorf("should not happen")
			}

			// the client reconnects to the server due to autoReconnector
			conn, err = ln.Accept()
			require.NoError(t, err)
			conn.Close()
		})
	}
}
