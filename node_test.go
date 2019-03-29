package gomavlib

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
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
		Dialect:          MustDialect([]Message{&MessageHeartbeat{}}),
		OutSystemId:      10,
		Endpoints:        []EndpointConf{t1},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:          MustDialect([]Message{&MessageHeartbeat{}}),
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
					if reflect.DeepEqual(e.Message(), testMsg1) == false ||
						e.SystemId() != 11 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
						return
					}
					node1.WriteMessageAll(testMsg2)
					step++

				case 1:
					if reflect.DeepEqual(e.Message(), testMsg3) == false ||
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
					if reflect.DeepEqual(e.Message(), testMsg2) == false ||
						e.SystemId() != 10 ||
						e.ComponentId() != 1 {
						t.Fatal("received wrong message")
						return
					}
					node2.WriteMessageAll(testMsg3)
					step++
				case 1:
					if reflect.DeepEqual(e.Message(), testMsg4) == false ||
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

type testLoopback chan []byte

func (ch testLoopback) Close() error {
	close(ch)
	return nil
}

func (ch testLoopback) Read(buf []byte) (int, error) {
	ret, ok := <-ch
	if ok == false {
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
		Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
		OutSystemId: 11,
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
	success := false
	func() {
		node1, err := NewNode(NodeConf{
			Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
			OutSystemId: 10,
			Endpoints: []EndpointConf{
				EndpointUdpServer{"127.0.0.1:5600"},
			},
			HeartbeatDisable: true,
		})
		require.NoError(t, err)
		defer node1.Close()

		node2, err := NewNode(NodeConf{
			Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
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

	if success == false {
		t.Fatalf("failed")
	}
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
		Dialect: MustDialect([]Message{&MessageHeartbeat{}}),
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InSignatureKey:   key2,
		OutSystemId:      10,
		OutSignatureKey:  key1,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect: MustDialect([]Message{&MessageHeartbeat{}}),
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		InSignatureKey:   key1,
		OutSystemId:      11,
		OutSignatureKey:  key2,
	})
	require.NoError(t, err)

	success := false
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer node1.Close()

		_, ok := <-node1.Events()
		if ok == false {
			return
		}
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		time.Sleep(500 * time.Millisecond)

		node2.WriteMessageAll(testMsg)

		_, ok := <-node2.Events()
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
		Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
		OutSystemId: 10,
		Endpoints: []EndpointConf{
			EndpointUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node2, err := NewNode(NodeConf{
		Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
		OutSystemId: 11,
		Endpoints: []EndpointConf{
			EndpointUdpServer{"127.0.0.1:5600"},
			EndpointUdpClient{"127.0.0.1:5601"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	node3, err := NewNode(NodeConf{
		Dialect:     MustDialect([]Message{&MessageHeartbeat{}}),
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
				node2.WriteFrameExcept(e.Channel(), e.Frame)
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

	if success == false {
		t.Fatal("failed")
	}
}
