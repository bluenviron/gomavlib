package gomavlib

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/pion/transport/v2/udp"
	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

type dummyReadWriter struct {
	chOut     chan []byte
	chIn      chan []byte
	chReadErr chan struct{}
}

func newDummyReadWriterPair() (*dummyReadWriter, *dummyReadWriter) {
	one := &dummyReadWriter{
		chOut:     make(chan []byte),
		chIn:      make(chan []byte),
		chReadErr: make(chan struct{}),
	}

	two := &dummyReadWriter{
		chOut:     one.chIn,
		chIn:      one.chOut,
		chReadErr: make(chan struct{}),
	}

	return one, two
}

func (e *dummyReadWriter) simulateReadError() {
	close(e.chReadErr)
}

func (e *dummyReadWriter) Close() error {
	close(e.chOut)
	close(e.chIn)
	return nil
}

func (e *dummyReadWriter) Read(p []byte) (int, error) {
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

func (e *dummyReadWriter) Write(p []byte) (int, error) {
	e.chIn <- p
	return len(p), nil
}

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

			go func() {
				conn, err := ln.Accept()
				require.NoError(t, err)
				defer conn.Close()

				dialectRW := &dialect.ReadWriter{Dialect: testDialect}
				err = dialectRW.Initialize()
				require.NoError(t, err)

				rw := &frame.ReadWriter{
					ByteReadWriter: conn,
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

				for i := 0; i < 3; i++ {
					var fr frame.Frame
					fr, err = rw.Read()
					require.NoError(t, err)
					require.Equal(t, &frame.V2Frame{
						SequenceNumber: byte(i),
						SystemID:       10,
						ComponentID:    1,
						Message: &MessageHeartbeat{
							Type:           1,
							Autopilot:      2,
							BaseMode:       3,
							CustomMode:     6,
							SystemStatus:   4,
							MavlinkVersion: 5,
						},
						Checksum: fr.GetChecksum(),
					}, fr)

					err = sw.Write(&MessageHeartbeat{
						Type:           6,
						Autopilot:      5,
						BaseMode:       4,
						CustomMode:     3,
						SystemStatus:   2,
						MavlinkVersion: 1,
					})
					require.NoError(t, err)
				}
			}()

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

			for i := 0; i < 3; i++ {
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
				require.Equal(t, &EventFrame{
					Frame: &frame.V2Frame{
						SequenceNumber: byte(i),
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
						Checksum: evt.(*EventFrame).Frame.GetChecksum(),
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

			closed := make(chan struct{})
			reconnected := make(chan struct{})

			go func() {
				conn, err2 := ln.Accept()
				require.NoError(t, err2)

				dialectRW := &dialect.ReadWriter{Dialect: testDialect}
				err2 = dialectRW.Initialize()
				require.NoError(t, err2)

				rw := &frame.Reader{
					ByteReader: conn,
					DialectRW:  dialectRW,
				}
				err2 = rw.Initialize()
				require.NoError(t, err2)

				fr, err2 := rw.Read()
				require.NoError(t, err2)
				require.Equal(t, &frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       10,
					ComponentID:    1,
					Message: &MessageHeartbeat{
						Type:           1,
						Autopilot:      2,
						BaseMode:       3,
						CustomMode:     6,
						SystemStatus:   4,
						MavlinkVersion: 5,
					},
					Checksum: fr.GetChecksum(),
				}, fr)

				_, err2 = rw.Read()
				require.Equal(t, io.EOF, err2)
				conn.Close()

				close(closed)

				// the client reconnects to the server due to autoReconnector
				conn, err2 = ln.Accept()
				require.NoError(t, err2)
				conn.Close()

				close(reconnected)
			}()

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

			err = node.WriteMessageAll(&MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     6,
				SystemStatus:   4,
				MavlinkVersion: 5,
			})
			require.NoError(t, err)

			select {
			case <-closed:
			case <-time.After(1 * time.Second):
				t.Errorf("should not happen")
			}

			<-node.Events()

			<-reconnected
		})
	}
}

func TestEndpointCustomClient(t *testing.T) {
	remote, local := newDummyReadWriterPair()

	node, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{EndpointCustomClient{
			Connect: func(_ context.Context) (net.Conn, error) {
				return remote, nil
			},
		}},
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
