package gomavlib

import (
	"bytes"
	"io"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/dialect"
	"github.com/aler9/gomavlib/msg"
)

type MAV_TYPE int
type MAV_AUTOPILOT int
type MAV_MODE_FLAG int
type MAV_STATE int

type MessageTest5 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest5) GetId() uint32 {
	return 5
}

type MessageTest6 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest6) GetId() uint32 {
	return 0x0607
}

type MessageTest8 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest8) GetId() uint32 {
	return 8
}

type MessageHeartbeat struct {
	Type           MAV_TYPE      `mavenum:"uint8"`
	Autopilot      MAV_AUTOPILOT `mavenum:"uint8"`
	BaseMode       MAV_MODE_FLAG `mavenum:"uint8"`
	CustomMode     uint32
	SystemStatus   MAV_STATE `mavenum:"uint8"`
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetId() uint32 {
	return 0
}

type MessageOpticalFlow struct {
	TimeUsec       uint64
	SensorId       uint8
	FlowX          int16
	FlowY          int16
	FlowCompMX     float32
	FlowCompMY     float32
	Quality        uint8
	GroundDistance float32
	FlowRateX      float32 `mavext:"true"`
	FlowRateY      float32 `mavext:"true"`
}

func (*MessageOpticalFlow) GetId() uint32 {
	return 100
}

type MessageRequestDataStream struct {
	TargetSystem    uint8
	TargetComponent uint8
	ReqStreamId     uint8
	ReqMessageRate  uint16
	StartStop       uint8
}

func (*MessageRequestDataStream) GetId() uint32 {
	return 66
}

var testDialect = &dialect.Dialect{3, []msg.Message{
	&MessageTest5{},
	&MessageTest6{},
	&MessageTest8{},
	&MessageHeartbeat{},
	&MessageOpticalFlow{},
}}

func doTest(t *testing.T, t1 EndpointConf, t2 EndpointConf) {
	var testMsg1 = &MessageHeartbeat{
		Type:           1,
		Autopilot:      2,
		BaseMode:       3,
		CustomMode:     6,
		SystemStatus:   4,
		MavlinkVersion: 5,
	}
	var testMsg2 = &MessageHeartbeat{
		Type:           6,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}
	var testMsg3 = &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}
	var testMsg4 = &MessageHeartbeat{
		Type:           7,
		Autopilot:      6,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect:          &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:       V2,
		OutSystemId:      10,
		Endpoints:        []EndpointConf{t1},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:          &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:       V2,
		OutSystemId:      11,
		Endpoints:        []EndpointConf{t2},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	success := false
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer node1.Close()

		step := 0

		for evt := range node1.Events() {
			switch e := evt.(type) {
			case *EventFrame:
				switch step {
				case 0:
					if !reflect.DeepEqual(e.Message(), testMsg1) ||
						e.SystemId() != 11 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
						return
					}
					node1.WriteMessageAll(testMsg2)
					step++

				case 1:
					if !reflect.DeepEqual(e.Message(), testMsg3) ||
						e.SystemId() != 11 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
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
			switch e := evt.(type) {
			case *EventFrame:
				switch step {
				case 0:
					if !reflect.DeepEqual(e.Message(), testMsg2) ||
						e.SystemId() != 10 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
						return
					}
					node2.WriteMessageAll(testMsg3)
					step++

				case 1:
					if !reflect.DeepEqual(e.Message(), testMsg4) ||
						e.SystemId() != 10 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
						return
					}
					success = true
					return
				}
			}
		}
	}()

	wg.Wait()

	require.Equal(t, true, success)
}

func TestNodeTcpServerClient(t *testing.T) {
	doTest(t, EndpointTcpServer{"127.0.0.1:5601"}, EndpointTcpClient{"127.0.0.1:5601"})
}

func TestNodeUdpServerClient(t *testing.T) {
	doTest(t, EndpointUdpServer{"127.0.0.1:5601"}, EndpointUdpClient{"127.0.0.1:5601"})
}

func TestNodeUdpBroadcastBroadcast(t *testing.T) {
	doTest(t, EndpointUdpBroadcast{"127.255.255.255:5602", ":5601"},
		EndpointUdpBroadcast{"127.255.255.255:5601", ":5602"})
}

type testLoopback chan []byte

func (ch testLoopback) Close() error {
	close(ch)
	return nil
}

func (ch testLoopback) Read(buf []byte) (int, error) {
	ret, ok := <-ch
	if !ok {
		return 0, errorTerminated
	}
	n := copy(buf, ret)
	return n, nil
}

func (ch testLoopback) Write(buf []byte) (int, error) {
	ch <- buf
	return len(buf), nil
}

type testEndpoint struct {
	io.ReadCloser
	io.Writer
}

func TestNodeCustomCustom(t *testing.T) {
	l1 := make(testLoopback)
	l2 := make(testLoopback)
	doTest(t, EndpointCustom{&testEndpoint{l1, l2}},
		EndpointCustom{&testEndpoint{l2, l1}})
}

func TestNodeError(t *testing.T) {
	_, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.Error(t, err)
}

func TestNodeCloseInLoop(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
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
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
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
	key1 := NewKey(bytes.Repeat([]byte("\x4F"), 32))
	key2 := NewKey(bytes.Repeat([]byte("\xA8"), 32))

	var testMsg = &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect: &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InKey:            key2,
		OutVersion:       V2,
		OutSystemId:      10,
		OutKey:           key1,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect: &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InKey:            key1,
		OutVersion:       V2,
		OutSystemId:      11,
		OutKey:           key2,
	})
	require.NoError(t, err)

	success := false
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer node1.Close()

		_, ok := <-node1.Events()
		if !ok {
			return
		}
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg)

		_, ok := <-node2.Events()
		if !ok {
			return
		}

		success = true
	}()

	wg.Wait()

	require.Equal(t, true, success)
}

func TestNodeRouting(t *testing.T) {
	var testMsg = &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 10,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
			EndpointUdpClient{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node3, err := NewNode(NodeConf{
		Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
		OutVersion:  V2,
		OutSystemId: 12,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	// wait client connection
	time.Sleep(500 * time.Millisecond)

	node1.WriteMessageAll(testMsg)

	success := false
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for evt := range node2.Events() {
			switch e := evt.(type) {
			case *EventFrame:
				node2.WriteFrameExcept(e.Channel, e.Frame)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		for evt := range node3.Events() {
			switch e := evt.(type) {
			case *EventFrame:
				if _, ok := e.Message().(*MessageHeartbeat); !ok ||
					e.SystemId() != 10 ||
					e.ComponentId() != 1 {
					t.Fatal("wrong message received")
				}
				success = true
				return
			}
		}
	}()

	wg.Wait()
	node1.Close()
	node2.Close()
	node3.Close()

	require.Equal(t, true, success)
}

func TestNodeHeartbeat(t *testing.T) {
	success := false

	func() {
		node1, err := NewNode(NodeConf{
			Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
			OutVersion:  V2,
			OutSystemId: 10,
			Endpoints: []EndpointConf{
				EndpointUdpServer{"127.0.0.1:5600"},
			},
			HeartbeatDisable: true,
		})
		require.NoError(t, err)
		defer node1.Close()

		node2, err := NewNode(NodeConf{
			Dialect:     &dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}},
			OutVersion:  V2,
			OutSystemId: 11,
			Endpoints: []EndpointConf{
				EndpointUdpClient{"127.0.0.1:5600"},
			},
			HeartbeatDisable: false,
			HeartbeatPeriod:  500 * time.Millisecond,
		})
		require.NoError(t, err)
		defer node2.Close()

		for evt := range node1.Events() {
			if ee, ok := evt.(*EventFrame); ok {
				if _, ok = ee.Message().(*MessageHeartbeat); ok {
					success = true
					break
				}
			}
		}
	}()

	require.Equal(t, true, success)
}

func TestNodeStreamRequest(t *testing.T) {
	success := false

	func() {
		node1, err := NewNode(NodeConf{
			Dialect: &dialect.Dialect{3, []msg.Message{
				&MessageHeartbeat{},
				&MessageRequestDataStream{},
			}},
			OutVersion:  V2,
			OutSystemId: 10,
			Endpoints: []EndpointConf{
				EndpointUdpServer{"127.0.0.1:5600"},
			},
			HeartbeatDisable:    true,
			StreamRequestEnable: true,
		})
		require.NoError(t, err)
		defer node1.Close()

		node2, err := NewNode(NodeConf{
			Dialect: &dialect.Dialect{3, []msg.Message{
				&MessageHeartbeat{},
				&MessageRequestDataStream{},
			}},
			OutVersion:  V2,
			OutSystemId: 10,
			Endpoints: []EndpointConf{
				EndpointUdpClient{"127.0.0.1:5600"},
			},
			HeartbeatDisable:       false,
			HeartbeatPeriod:        500 * time.Millisecond,
			HeartbeatAutopilotType: 3, // MAV_AUTOPILOT_ARDUPILOTMEGA
		})
		require.NoError(t, err)
		defer node2.Close()

		go func() {
			for range node1.Events() {
			}
		}()

		for evt := range node2.Events() {
			if ee, ok := evt.(*EventFrame); ok {
				if _, ok = ee.Message().(*MessageRequestDataStream); ok {
					success = true
					break
				}
			}
		}
	}()

	require.Equal(t, true, success)
}
