package gomavlib

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

func TestEndpointSerial(t *testing.T) {
	done := make(chan struct{})
	n := 0

	serialOpenFunc = func(_ string, _ int) (io.ReadWriteCloser, error) {
		remote, local := newDummyReadWriterPair()

		n++
		switch n {
		case 0:
			return remote, nil

		case 1:

		default:
			t.Errorf("should not happen")
		}

		go func() {
			dialectRW := &dialect.ReadWriter{Dialect: testDialect}
			err := dialectRW.Initialize()
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

			for i := range 3 {
				err = sw.Write(&MessageHeartbeat{
					Type:           1,
					Autopilot:      2,
					BaseMode:       3,
					CustomMode:     6,
					SystemStatus:   4,
					MavlinkVersion: 5,
				})
				require.NoError(t, err)

				var fr frame.Frame
				fr, err = rw.Read()
				require.NoError(t, err)
				require.Equal(t, &frame.V2Frame{
					SequenceNumber: byte(i),
					SystemID:       10,
					ComponentID:    1,
					Message: &MessageHeartbeat{
						Type:           6,
						Autopilot:      5,
						BaseMode:       4,
						CustomMode:     3,
						SystemStatus:   2,
						MavlinkVersion: 1,
					},
					Checksum: fr.GetChecksum(),
				}, fr)
			}

			close(done)
		}()

		return remote, nil
	}

	node, err := NewNode(NodeConf{
		Dialect:     testDialect,
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

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	for i := range 3 {
		evt = <-node.Events()
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceNumber: byte(i),
				SystemID:       11,
				ComponentID:    1,
				Message: &MessageHeartbeat{
					Type:           1,
					Autopilot:      2,
					BaseMode:       3,
					CustomMode:     6,
					SystemStatus:   4,
					MavlinkVersion: 5,
				},
				Checksum: evt.(*EventFrame).Frame.GetChecksum(),
			},
			Channel: evt.(*EventFrame).Channel,
		}, evt)

		err = node.WriteMessageAll(&MessageHeartbeat{
			Type:           6,
			Autopilot:      5,
			BaseMode:       4,
			CustomMode:     3,
			SystemStatus:   2,
			MavlinkVersion: 1,
		})
		require.NoError(t, err)
	}

	<-done
}

func TestEndpointSerialReconnect(t *testing.T) {
	done := make(chan struct{})
	count := 0

	serialOpenFunc = func(_ string, _ int) (io.ReadWriteCloser, error) {
		remote, local := newDummyReadWriterPair()

		switch count {
		case 0:
			go func() {
				dialectRW := &dialect.ReadWriter{Dialect: testDialect}
				err := dialectRW.Initialize()
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

				err = sw.Write(&MessageHeartbeat{
					Type:           1,
					Autopilot:      2,
					BaseMode:       3,
					CustomMode:     6,
					SystemStatus:   4,
					MavlinkVersion: 5,
				})
				require.NoError(t, err)

				fr, err := rw.Read()
				require.NoError(t, err)
				require.Equal(t, &frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       10,
					ComponentID:    1,
					Message: &MessageHeartbeat{
						Type:           6,
						Autopilot:      5,
						BaseMode:       4,
						CustomMode:     3,
						SystemStatus:   2,
						MavlinkVersion: 1,
					},
					Checksum: fr.GetChecksum(),
				}, fr)

				remote.simulateReadError()
			}()

		case 1:
			go func() {
				dialectRW := &dialect.ReadWriter{Dialect: testDialect}
				err := dialectRW.Initialize()
				require.NoError(t, err)

				rw := &frame.Reader{
					ByteReader: local,
					DialectRW:  dialectRW,
				}
				err = rw.Initialize()
				require.NoError(t, err)

				fr, err := rw.Read()
				require.NoError(t, err)
				require.Equal(t, &frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       10,
					ComponentID:    1,
					Message: &MessageHeartbeat{
						Type:           7,
						Autopilot:      5,
						BaseMode:       4,
						CustomMode:     3,
						SystemStatus:   2,
						MavlinkVersion: 1,
					},
					Checksum: fr.GetChecksum(),
				}, fr)

				close(done)
			}()
		}

		count++
		return remote, nil
	}

	node, err := NewNode(NodeConf{
		Dialect:     testDialect,
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

	evt := <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	evt = <-node.Events()
	require.Equal(t, &EventFrame{
		Frame: &frame.V2Frame{
			SequenceNumber: 0,
			SystemID:       11,
			ComponentID:    1,
			Message: &MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     6,
				SystemStatus:   4,
				MavlinkVersion: 5,
			},
			Checksum: evt.(*EventFrame).Frame.GetChecksum(),
		},
		Channel: evt.(*EventFrame).Channel,
	}, evt)

	err = node.WriteMessageAll(&MessageHeartbeat{
		Type:           6,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	})
	require.NoError(t, err)

	evt = <-node.Events()
	require.Equal(t, &EventChannelClose{
		Channel: evt.(*EventChannelClose).Channel,
		Error:   evt.(*EventChannelClose).Error,
	}, evt)

	evt = <-node.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	err = node.WriteMessageAll(&MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	})
	require.NoError(t, err)

	<-done
}
