package gomavlib

import (
	"bytes"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

type (
	MAV_TYPE      uint64 //nolint:revive
	MAV_AUTOPILOT uint64 //nolint:revive
	MAV_MODE_FLAG uint64 //nolint:revive
	MAV_STATE     uint64 //nolint:revive
)

type MessageHeartbeat struct {
	Type           MAV_TYPE      `mavenum:"uint8"`
	Autopilot      MAV_AUTOPILOT `mavenum:"uint8"`
	BaseMode       MAV_MODE_FLAG `mavenum:"uint8"`
	CustomMode     uint32
	SystemStatus   MAV_STATE `mavenum:"uint8"`
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetID() uint32 {
	return 0
}

type MessageRequestDataStream struct {
	TargetSystem    uint8
	TargetComponent uint8
	ReqStreamId     uint8 //nolint:revive
	ReqMessageRate  uint16
	StartStop       uint8
}

func (*MessageRequestDataStream) GetID() uint32 {
	return 66
}

var testDialect = &dialect.Dialect{
	Version:  3,
	Messages: []message.Message{&MessageHeartbeat{}},
}

var testMessage = &MessageHeartbeat{
	Type:           7,
	Autopilot:      5,
	BaseMode:       4,
	CustomMode:     3,
	SystemStatus:   2,
	MavlinkVersion: 1,
}

func TestNodeError(t *testing.T) {
	_, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.Error(t, err)
}

func TestNodeCloseInLoop(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	err = node2.WriteMessageAll(testMessage)
	require.NoError(t, err)

	for evt := range node1.Events() {
		if _, ok := evt.(*EventChannelOpen); ok {
			node1.Close()
		}
	}
}

func TestNodeWriteAll(t *testing.T) {
	for _, ca := range []string{"message", "frame"} {
		t.Run(ca, func(t *testing.T) {
			server, err := NewNode(NodeConf{
				Dialect:     testDialect,
				OutVersion:  V2,
				OutSystemID: 11,
				Endpoints: []EndpointConf{
					EndpointTCPServer{"127.0.0.1:5600"},
				},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer server.Close()

			var wg sync.WaitGroup
			wg.Add(5)

			for range 5 {
				var client *Node
				client, err = NewNode(NodeConf{
					Dialect:     testDialect,
					OutVersion:  V2,
					OutSystemID: 11,
					Endpoints: []EndpointConf{
						EndpointTCPClient{"127.0.0.1:5600"},
					},
					HeartbeatDisable: true,
				})
				require.NoError(t, err)
				defer client.Close()

				go func() {
					for evt := range client.Events() {
						if fr, ok := evt.(*EventFrame); ok {
							require.Equal(t, &EventFrame{
								Frame: &frame.V2Frame{
									SequenceNumber: 0,
									SystemID:       11,
									ComponentID:    1,
									Message:        testMessage,
									Checksum:       fr.Frame.GetChecksum(),
								},
								Channel: fr.Channel,
							}, fr)
							wg.Done()
						}
					}
				}()
			}

			count := 0
			for evt := range server.Events() {
				if _, ok := evt.(*EventChannelOpen); ok {
					count++
					if count == 5 {
						break
					}
				}
			}

			if ca == "message" {
				err = server.WriteMessageAll(testMessage)
				require.NoError(t, err)
			} else {
				err = server.WriteFrameAll(&frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       11,
					ComponentID:    1,
					Message:        testMessage,
					Checksum:       55967,
				})
				require.NoError(t, err)
			}
			wg.Wait()
		})
	}
}

func TestNodeWriteExcept(t *testing.T) {
	for _, ca := range []string{"message", "frame"} {
		t.Run(ca, func(t *testing.T) {
			server, err := NewNode(NodeConf{
				Dialect:     testDialect,
				OutVersion:  V2,
				OutSystemID: 11,
				Endpoints: []EndpointConf{
					EndpointTCPServer{"127.0.0.1:5600"},
				},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer server.Close()

			var wg sync.WaitGroup
			wg.Add(4)

			for range 5 {
				var client *Node
				client, err = NewNode(NodeConf{
					Dialect:     testDialect,
					OutVersion:  V2,
					OutSystemID: 11,
					Endpoints: []EndpointConf{
						EndpointTCPClient{"127.0.0.1:5600"},
					},
					HeartbeatDisable: true,
				})
				require.NoError(t, err)
				defer client.Close()

				go func() {
					for evt := range client.Events() {
						if fr, ok := evt.(*EventFrame); ok {
							require.Equal(t, &EventFrame{
								Frame: &frame.V2Frame{
									SequenceNumber: 0,
									SystemID:       11,
									ComponentID:    1,
									Message:        testMessage,
									Checksum:       fr.Frame.GetChecksum(),
								},
								Channel: fr.Channel,
							}, fr)
							wg.Done()
						}
					}
				}()
			}

			count := 0
			var except *Channel
			for evt := range server.Events() {
				if evt2, ok := evt.(*EventChannelOpen); ok {
					if count == 1 {
						except = evt2.Channel
					}
					count++
					if count == 5 {
						break
					}
				}
			}

			if ca == "message" {
				err = server.WriteMessageExcept(except, testMessage)
				require.NoError(t, err)
			} else {
				err = server.WriteFrameExcept(except, &frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       11,
					ComponentID:    1,
					Message:        testMessage,
					Checksum:       55967,
				})
				require.NoError(t, err)
			}
			wg.Wait()
		})
	}
}

func TestNodeWriteTo(t *testing.T) {
	for _, ca := range []string{"message", "frame"} {
		t.Run(ca, func(t *testing.T) {
			server, err := NewNode(NodeConf{
				Dialect:     testDialect,
				OutVersion:  V2,
				OutSystemID: 11,
				Endpoints: []EndpointConf{
					EndpointTCPServer{"127.0.0.1:5600"},
				},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer server.Close()

			recv := make(chan struct{})

			for range 5 {
				var client *Node
				client, err = NewNode(NodeConf{
					Dialect:     testDialect,
					OutVersion:  V2,
					OutSystemID: 11,
					Endpoints: []EndpointConf{
						EndpointTCPClient{"127.0.0.1:5600"},
					},
					HeartbeatDisable: true,
				})
				require.NoError(t, err)
				defer client.Close()

				go func() {
					for evt := range client.Events() {
						if fr, ok := evt.(*EventFrame); ok {
							require.Equal(t, &EventFrame{
								Frame: &frame.V2Frame{
									SequenceNumber: 0,
									SystemID:       11,
									ComponentID:    1,
									Message:        testMessage,
									Checksum:       fr.Frame.GetChecksum(),
								},
								Channel: fr.Channel,
							}, fr)
							close(recv)
						}
					}
				}()
			}

			count := 0
			var except *Channel
			for evt := range server.Events() {
				if evt2, ok := evt.(*EventChannelOpen); ok {
					if count == 1 {
						except = evt2.Channel
					}
					count++
					if count == 5 {
						break
					}
				}
			}

			if ca == "message" {
				err = server.WriteMessageTo(except, testMessage)
				require.NoError(t, err)
			} else {
				err = server.WriteFrameTo(except, &frame.V2Frame{
					SequenceNumber: 0,
					SystemID:       11,
					ComponentID:    1,
					Message:        testMessage,
					Checksum:       55967,
				})
				require.NoError(t, err)
			}
			<-recv
		})
	}
}

func TestNodeWriteMessageInLoop(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	err = node2.WriteMessageAll(testMessage)
	require.NoError(t, err)

	for evt := range node1.Events() {
		if _, ok := evt.(*EventChannelOpen); ok {
			for range 10 {
				err = node1.WriteMessageAll(testMessage)
				require.NoError(t, err)
			}
			break
		}
	}
}

func TestNodeSignature(t *testing.T) {
	key1 := frame.NewV2Key(bytes.Repeat([]byte("\x4F"), 32))
	key2 := frame.NewV2Key(bytes.Repeat([]byte("\xA8"), 32))

	node1, err := NewNode(NodeConf{
		Dialect: testDialect,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InKey:            key2,
		OutVersion:       V2,
		OutSystemID:      10,
		OutKey:           key1,
	})
	require.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(NodeConf{
		Dialect: testDialect,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InKey:            key1,
		OutVersion:       V2,
		OutSystemID:      11,
		OutKey:           key2,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	err = node2.WriteMessageAll(testMessage)
	require.NoError(t, err)

	<-node1.Events()
	evt = <-node1.Events()
	fr, ok := evt.(*EventFrame)
	require.Equal(t, true, ok)
	require.Equal(t, &EventFrame{
		Frame: &frame.V2Frame{
			SequenceNumber:      0,
			SystemID:            11,
			ComponentID:         1,
			Message:             testMessage,
			Checksum:            fr.Frame.GetChecksum(),
			IncompatibilityFlag: 1,
			SignatureLinkID:     fr.Frame.(*frame.V2Frame).SignatureLinkID,
			Signature:           fr.Frame.(*frame.V2Frame).Signature,
			SignatureTimestamp:  fr.Frame.(*frame.V2Frame).SignatureTimestamp,
		},
		Channel: fr.Channel,
	}, evt)
}

func TestNodeRoute(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
			EndpointUDPClient{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	node3, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 12,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node3.Close()

	evt := <-node1.Events()
	_, ok := evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	err = node1.WriteMessageAll(testMessage)
	require.NoError(t, err)

	evt = <-node2.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-node2.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-node2.Events()
	fr, ok := evt.(*EventFrame)
	require.Equal(t, true, ok)

	err = node2.WriteFrameExcept(fr.Channel, fr.Frame)
	require.NoError(t, err)

	evt = <-node3.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-node3.Events()
	fr, ok = evt.(*EventFrame)
	require.Equal(t, true, ok)
	require.Equal(t, &frame.V2Frame{
		SystemID:    10,
		ComponentID: 1,
		Message:     testMessage,
		Checksum:    fr.Frame.GetChecksum(),
	}, fr.Frame)
}

func TestNodeFixFrame(t *testing.T) {
	key := frame.NewV2Key(bytes.Repeat([]byte("\xB2"), 32))

	node1, err := NewNode(NodeConf{
		Dialect: testDialect,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InKey:            key,
		OutVersion:       V2,
		OutSystemID:      10,
	})
	require.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(NodeConf{
		Dialect: testDialect,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		OutVersion:       V2,
		OutSystemID:      11,
		OutKey:           key,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	require.Equal(t, &EventChannelOpen{
		Channel: evt.(*EventChannelOpen).Channel,
	}, evt)

	fra := &frame.V2Frame{
		SequenceNumber:      13,
		SystemID:            15,
		ComponentID:         11,
		Message:             testMessage,
		IncompatibilityFlag: 1,
		SignatureLinkID:     16,
		SignatureTimestamp:  23434542,
		Signature:           (*frame.V2Signature)(&[6]byte{0, 0, 0, 0, 0, 0}),
	}

	err = node2.FixFrame(fra)
	require.NoError(t, err)

	err = node2.WriteFrameAll(fra)
	require.NoError(t, err)

	<-node1.Events()
	evt = <-node1.Events()
	fr, ok := evt.(*EventFrame)
	require.Equal(t, true, ok)
	require.Equal(t, &EventFrame{
		Frame: &frame.V2Frame{
			SequenceNumber:      13,
			SystemID:            15,
			ComponentID:         11,
			Message:             testMessage,
			Checksum:            fr.Frame.GetChecksum(),
			IncompatibilityFlag: 1,
			SignatureLinkID:     16,
			Signature:           fr.Frame.(*frame.V2Frame).Signature,
			SignatureTimestamp:  23434542,
		},
		Channel: fr.Channel,
	}, evt)
}

func TestNodeWriteSameToMultiple(t *testing.T) {
	server, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointTCPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer server.Close()

	client1, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointTCPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer client1.Close()

	client2, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointTCPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer client2.Close()

	evt := <-client1.Events()
	_, ok := evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-client2.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-server.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	evt = <-server.Events()
	_, ok = evt.(*EventChannelOpen)
	require.Equal(t, true, ok)

	fr := &frame.V2Frame{
		SequenceNumber: 0,
		SystemID:       11,
		ComponentID:    1,
		Message:        testMessage,
		Checksum:       55967,
	}

	err = client1.WriteFrameAll(fr)
	require.NoError(t, err)

	err = client2.WriteFrameAll(fr)
	require.NoError(t, err)

	for range 2 {
		evt = <-server.Events()
		var fr *EventFrame
		fr, ok = evt.(*EventFrame)
		require.Equal(t, true, ok)
		require.Equal(t, &EventFrame{
			Frame: &frame.V2Frame{
				SequenceNumber: 0,
				SystemID:       11,
				ComponentID:    1,
				Message:        testMessage,
				Checksum:       55967,
			},
			Channel: fr.Channel,
		}, evt)
	}
}
