package gomavlib

import (
	"bytes"
	"reflect"
	"sync"
	"testing"
	"time"
)

func doTest(t *testing.T, t1 TransportConf, t2 TransportConf) {
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
		Transports:       []TransportConf{t1},
		HeartbeatDisable: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	node2, err := NewNode(NodeConf{
		Dialect:          []Message{&MessageHeartbeat{}},
		SystemId:         11,
		ComponentId:      1,
		Transports:       []TransportConf{t2},
		HeartbeatDisable: true,
	})
	if err != nil {
		t.Fatal(err)
	}

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

		node1.WriteMessage(nil, testMsg2)

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

		node1.WriteMessage(nil, testMsg4)
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		// wait connection to server
		time.Sleep(500 * time.Millisecond)

		node2.WriteMessage(nil, testMsg1)

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

		node2.WriteMessage(nil, testMsg3)

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
	doTest(t, TransportTcpServer{"127.0.0.1:5601"}, TransportTcpClient{"127.0.0.1:5601"})
}

func TestNodeUdpServerClient(t *testing.T) {
	doTest(t, TransportUdpServer{"127.0.0.1:5601"}, TransportUdpClient{"127.0.0.1:5601"})
}

func TestNodeUdpBroadcastBroadcast(t *testing.T) {
	doTest(t, TransportUdpBroadcast{"127.255.255.255:5602", ":5601"},
		TransportUdpBroadcast{"127.255.255.255:5601", ":5602"})
}

type CustomTransport struct {
	writeBuf *bytes.Buffer
	readChan chan struct{}
}

func NewCustomTransport() *CustomTransport {
	return &CustomTransport{
		readChan: make(chan struct{}),
		writeBuf: bytes.NewBuffer(nil),
	}
}

func (c *CustomTransport) Close() error {
	close(c.readChan)
	return nil
}

func (c *CustomTransport) Read(buf []byte) (int, error) {
	<-c.readChan
	return 0, errorTerminated
}

func (c *CustomTransport) Write(buf []byte) (int, error) {
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

	rwc := NewCustomTransport()

	func() {
		node, err := NewNode(NodeConf{
			Dialect:     []Message{&MessageHeartbeat{}},
			SystemId:    11,
			ComponentId: 1,
			Transports: []TransportConf{
				TransportCustom{rwc},
			},
			HeartbeatDisable: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		defer node.Close()

		node.WriteMessage(nil, testMsg)
	}()

	if reflect.DeepEqual(rwc.writeBuf.Bytes(), res) == false {
		t.Fatal("result different")
	}
}

func TestNodeError(t *testing.T) {
	_, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    11,
		ComponentId: 1,
		Transports: []TransportConf{
			TransportUdpServer{"127.0.0.1:5600"},
			TransportUdpServer{"127.0.0.1:5600"},
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
			Transports: []TransportConf{
				TransportUdpServer{"127.0.0.1:5600"},
			},
			HeartbeatDisable: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		defer node1.Close()

		node2, err := NewNode(NodeConf{
			Dialect:     []Message{&MessageHeartbeat{}},
			SystemId:    11,
			ComponentId: 1,
			Transports: []TransportConf{
				TransportUdpClient{"127.0.0.1:5600"},
			},
			HeartbeatDisable: false,
			HeartbeatPeriod:  500 * time.Millisecond,
		})
		if err != nil {
			t.Fatal(err)
		}
		defer node2.Close()

		_, ok := node1.Read()
		if ok == false {
			t.Fatal(err)
		}
	}()
}

func TestNodeSignature(t *testing.T) {
	key1 := NewSignatureKey(bytes.Repeat([]byte("\x4F"), 32))
	key2 := NewSignatureKey(bytes.Repeat([]byte("\xA8"), 32))

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
		Transports: []TransportConf{
			TransportUdpServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		SignatureInKey:   key2,
		SignatureOutKey:  key1,
	})
	if err != nil {
		t.Fatal(err)
	}

	node2, err := NewNode(NodeConf{
		Dialect:     []Message{&MessageHeartbeat{}},
		SystemId:    11,
		ComponentId: 1,
		Transports: []TransportConf{
			TransportUdpClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
		SignatureInKey:   key1,
		SignatureOutKey:  key2,
	})
	if err != nil {
		t.Fatal(err)
	}

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

		node1.WriteMessage(nil, testMsg)
	}()

	go func() {
		defer wg.Done()
		defer node2.Close()

		time.Sleep(500 * time.Millisecond)

		node2.WriteMessage(nil, testMsg)

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
