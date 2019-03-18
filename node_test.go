package gomavlib

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"reflect"
	"sync"
	"testing"
	"time"
)

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
		Dialect:          []Message{&MessageHeartbeat{}},
		SystemId:         10,
		ComponentId:      1,
		Endpoints:        []EndpointConf{t1},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:          []Message{&MessageHeartbeat{}},
		SystemId:         11,
		ComponentId:      1,
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

		res1, ok := node1.Read()
		if ok == false {
			return
		}
		if reflect.DeepEqual(res1.Message(), testMsg1) == false ||
			res1.SystemId() != 11 ||
			res1.ComponentId() != 1 {
			t.Fatal("received wrong message")
			return
		}

		node1.WriteMessageAll(testMsg2)

		res2, ok := node1.Read()
		if ok == false {
			return
		}
		if reflect.DeepEqual(res2.Message(), testMsg3) == false ||
			res2.SystemId() != 11 ||
			res2.ComponentId() != 1 {
			t.Fatal("received wrong message")
			return
		}

		if res1.Channel() != res2.Channel() {
			t.Fatal("message received on two different channels")
			return
		}

		node1.WriteMessageAll(testMsg4)
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		// wait connection to server
		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg1)

		res, ok := node2.Read()
		if ok == false {
			return
		}
		if reflect.DeepEqual(res.Message(), testMsg2) == false ||
			res.SystemId() != 10 ||
			res.ComponentId() != 1 {
			t.Fatal("received wrong message")
			return
		}

		node2.WriteMessageAll(testMsg3)

		res, ok = node2.Read()
		if ok == false {
			return
		}

		if reflect.DeepEqual(res.Message(), testMsg4) == false ||
			res.SystemId() != 10 ||
			res.ComponentId() != 1 {
			t.Fatal("received wrong message")
			return
		}

		success = true
	}()

	wg.Wait()

	if success == false {
		t.Fatalf("test failed")
	}
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

type CustomEndpoint struct {
	writeBuf *bytes.Buffer
	readChan chan struct{}
}

func NewCustomEndpoint() *CustomEndpoint {
	return &CustomEndpoint{
		readChan: make(chan struct{}),
		writeBuf: bytes.NewBuffer(nil),
	}
}

func (c *CustomEndpoint) Close() error {
	close(c.readChan)
	return nil
}

func (c *CustomEndpoint) Read(buf []byte) (int, error) {
	<-c.readChan
	return 0, errorTerminated
}

func (c *CustomEndpoint) Write(buf []byte) (int, error) {
	c.writeBuf.Write(buf)
	return len(buf), nil
}

func TestNodeCustom(t *testing.T) {
	var testMsg = &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}
	var res = []byte{253, 9, 0, 0, 0, 11, 1, 0, 0, 0, 3, 0, 0, 0, 7, 5, 4, 2, 1, 159, 218}

	rwc := NewCustomEndpoint()

	func() {
		node, err := NewNode(NodeConf{
			Dialect:     []Message{&MessageHeartbeat{}},
			SystemId:    11,
			ComponentId: 1,
			Endpoints: []EndpointConf{
				EndpointCustom{rwc},
			},
			HeartbeatDisable: true,
		})
		require.NoError(t, err)
		defer node.Close()

		node.WriteMessageAll(testMsg)
	}()

	require.Equal(t, res, rwc.writeBuf.Bytes())
}

func TestNodeError(t *testing.T) {
	_, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    11,
		ComponentId: 1,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})

	if err == nil {
		t.Fatal("no error")
	}
}

func TestNodeHeartbeat(t *testing.T) {
	func() {
		node1, err := NewNode(NodeConf{
			Dialect:     []Message{&MessageHeartbeat{}},
			SystemId:    10,
			ComponentId: 1,
			Endpoints: []EndpointConf{
				EndpointUdpServer{"127.0.0.1:5600"},
			},
			HeartbeatDisable: true,
		})
		require.NoError(t, err)
		defer node1.Close()

		node2, err := NewNode(NodeConf{
			Dialect:     []Message{&MessageHeartbeat{}},
			SystemId:    11,
			ComponentId: 1,
			Endpoints: []EndpointConf{
				EndpointUdpClient{"127.0.0.1:5600"},
			},
			HeartbeatDisable: false,
			HeartbeatPeriod:  500 * time.Millisecond,
		})
		require.NoError(t, err)
		defer node2.Close()

		_, ok := node1.Read()
		if ok == false {
			t.Fatal(err)
		}
	}()
}

func TestNodeFrameSignature(t *testing.T) {
	key1 := NewFrameSignatureKey(bytes.Repeat([]byte("\x4F"), 32))
	key2 := NewFrameSignatureKey(bytes.Repeat([]byte("\xA8"), 32))

	var testMsg = &MessageHeartbeat{
		Type:           7,
		Autopilot:      5,
		BaseMode:       4,
		CustomMode:     3,
		SystemStatus:   2,
		MavlinkVersion: 1,
	}

	node1, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    10,
		ComponentId: 1,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		SignatureInKey:   key2,
		SignatureOutKey:  key1,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    11,
		ComponentId: 1,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		SignatureInKey:   key1,
		SignatureOutKey:  key2,
	})
	require.NoError(t, err)

	success := false
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer node1.Close()

		_, ok := node1.Read()
		if ok == false {
			return
		}

		node1.WriteMessageAll(testMsg)
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg)

		_, ok := node2.Read()
		if ok == false {
			return
		}

		success = true
	}()

	wg.Wait()

	if success == false {
		t.Fatalf("test failed")
	}
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
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    10,
		ComponentId: 1,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    11,
		ComponentId: 1,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
			EndpointUdpClient{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node3, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    12,
		ComponentId: 1,
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

		res, ok := node2.Read()
		if ok == false {
			return
		}
		node2.WriteFrameExcept(res.Channel(), res.Frame())
	}()

	go func() {
		defer wg.Done()

		res, ok := node3.Read()
		if ok == false {
			return
		}

		if _, ok := res.Message().(*MessageHeartbeat); !ok ||
			res.SystemId() != 10 ||
			res.ComponentId() != 1 {
			t.Fatal("wrong message received")
		}

		success = true
	}()

	wg.Wait()
	node1.Close()
	node2.Close()
	node3.Close()

	if success == false {
		t.Fatal("failed")
	}
}
