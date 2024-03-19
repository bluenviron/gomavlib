package gomavlib

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
)

var _ endpointChannelProvider = (*endpointSerial)(nil)

func TestEndpointSerial(t *testing.T) {
	done := make(chan struct{})
	n := 0

	serialOpenFunc = func(_ string, _ int) (io.ReadWriteCloser, error) {
		remote, local := newDummyReadWriterPair()

		n++
		switch n {
		case 1:
			return remote, nil

		case 2:

		default:
			t.Errorf("should not happen")
		}

		go func() {
			dialectRW, err := dialect.NewReadWriter(testDialect)
			require.NoError(t, err)

			rw, err := frame.NewReadWriter(frame.ReadWriterConf{
				ReadWriter:  local,
				DialectRW:   dialectRW,
				OutVersion:  frame.V2,
				OutSystemID: 11,
			})
			require.NoError(t, err)

			for i := 0; i < 3; i++ {
				err = rw.WriteMessage(&MessageHeartbeat{
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
					SequenceID:  byte(i),
					SystemID:    10,
					ComponentID: 1,
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

	for i := 0; i < 3; i++ {
		evt := <-node.Events()
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceID:  byte(i),
				SystemID:    11,
				ComponentID: 1,
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

		err := node.WriteMessageAll(&MessageHeartbeat{
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
		case 0: // skip first call to serialOpenFunc()

		case 1:
			go func() {
				dialectRW, err := dialect.NewReadWriter(testDialect)
				require.NoError(t, err)

				rw, err := frame.NewReadWriter(frame.ReadWriterConf{
					ReadWriter:  local,
					DialectRW:   dialectRW,
					OutVersion:  frame.V2,
					OutSystemID: 11,
				})
				require.NoError(t, err)

				err = rw.WriteMessage(&MessageHeartbeat{
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
					SequenceID:  0,
					SystemID:    10,
					ComponentID: 1,
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

		case 2:
			go func() {
				dialectRW, err := dialect.NewReadWriter(testDialect)
				require.NoError(t, err)

				rw, err := frame.NewReadWriter(frame.ReadWriterConf{
					ReadWriter:  local,
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
			SequenceID:  0,
			SystemID:    11,
			ComponentID: 1,
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
