package gomavlib

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/message"
)

type (
	MAV_TYPE      uint32 //nolint:revive
	MAV_AUTOPILOT uint32 //nolint:revive
	MAV_MODE_FLAG uint32 //nolint:revive
	MAV_STATE     uint32 //nolint:revive
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

func doTest(t *testing.T, t1 EndpointConf, t2 EndpointConf) {
	testMsg1 := &MessageHeartbeat{
		Type:           1,
		Autopilot:      2,
		BaseMode:       3,
		CustomMode:     6,
		SystemStatus:   4,
		MavlinkVersion: 5,
	}
	testMsg2 := &MessageHeartbeat{
		Type:           6,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}
	testMsg3 := &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}
	testMsg4 := &MessageHeartbeat{
		Type:           7,
		Autopilot:      6,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect:          &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{t1},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:          &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:       V2,
		OutSystemID:      11, //nolint:revive
		Endpoints:        []EndpointConf{t2},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2)
	var err1 error
	var err2 error

	go func() {
		defer wg.Done()
		defer node1.Close()

		step := 0

		for evt := range node1.Events() {
			if e, ok := evt.(*EventFrame); ok {
				switch step {
				case 0:
					if !reflect.DeepEqual(e.Message(), testMsg1) ||
						e.SystemID() != 11 ||
						e.ComponentID() != 1 {
						err1 = fmt.Errorf("received wrong message")
						return
					}
					node1.WriteMessageAll(testMsg2)
					step++

				case 1:
					if !reflect.DeepEqual(e.Message(), testMsg3) ||
						e.SystemID() != 11 ||
						e.ComponentID() != 1 {
						err1 = fmt.Errorf("received wrong message")
						return
					}
					node1.WriteMessageAll(testMsg4)
					return
				}
			}
		}
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		// wait connection to server
		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg1)

		step := 0

		for evt := range node2.Events() {
			if e, ok := evt.(*EventFrame); ok {
				switch step {
				case 0:
					if !reflect.DeepEqual(e.Message(), testMsg2) ||
						e.SystemID() != 10 ||
						e.ComponentID() != 1 {
						err2 = fmt.Errorf("received wrong message")
						return
					}
					node2.WriteMessageAll(testMsg3)
					step++

				case 1:
					if !reflect.DeepEqual(e.Message(), testMsg4) ||
						e.SystemID() != 10 ||
						e.ComponentID() != 1 {
						err2 = fmt.Errorf("received wrong message")
						return
					}
					return
				}
			}
		}
	}()

	wg.Wait()
	require.NoError(t, err1)
	require.NoError(t, err2)
}

func TestNodeError(t *testing.T) {
	_, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
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
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	go func() {
		defer node2.Close()

		// wait connection to server
		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(&MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     6,
			SystemStatus:   4,
			MavlinkVersion: 5,
		})
	}()

	for evt := range node1.Events() {
		if _, ok := evt.(*EventChannelOpen); ok {
			node1.Close()
		}
	}
}

func TestNodeWriteMultipleInLoop(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	go func() {
		defer node2.Close()

		// wait connection to server
		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(&MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     6,
			SystemStatus:   4,
			MavlinkVersion: 5,
		})
	}()

	for evt := range node1.Events() {
		if _, ok := evt.(*EventChannelOpen); ok {
			for i := 0; i < 100; i++ {
				node1.WriteMessageAll(&MessageHeartbeat{
					Type:           1,
					Autopilot:      2,
					BaseMode:       3,
					CustomMode:     6,
					SystemStatus:   4,
					MavlinkVersion: 5,
				})
			}
			node1.Close()
		}
	}
}

func TestNodeSignature(t *testing.T) {
	key1 := frame.NewV2Key(bytes.Repeat([]byte("\x4F"), 32))
	key2 := frame.NewV2Key(bytes.Repeat([]byte("\xA8"), 32))

	testMsg := &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect: &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
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

	node2, err := NewNode(NodeConf{
		Dialect: &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer node1.Close()
		<-node1.Events()
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg)

		<-node2.Events()
	}()

	wg.Wait()
}

func TestNodeRouting(t *testing.T) {
	testMsg := &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
			EndpointUDPClient{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node3, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []message.Message{&MessageHeartbeat{}}}, //nolint:govet
		OutVersion:  V2,
		OutSystemID: 12,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	// wait client connection
	time.Sleep(500 * time.Millisecond)

	node1.WriteMessageAll(testMsg)

	var wg sync.WaitGroup
	wg.Add(2)
	var err2 error

	go func() {
		defer wg.Done()

		for evt := range node2.Events() {
			if e, ok := evt.(*EventFrame); ok {
				node2.WriteFrameExcept(e.Channel, e.Frame)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		for evt := range node3.Events() {
			if e, ok := evt.(*EventFrame); ok {
				if _, ok := e.Message().(*MessageHeartbeat); !ok ||
					e.SystemID() != 10 ||
					e.ComponentID() != 1 {
					err2 = fmt.Errorf("wrong message received")
					return
				}
				return
			}
		}
	}()

	wg.Wait()
	node1.Close()
	node2.Close()
	node3.Close()

	require.NoError(t, err2)
}
