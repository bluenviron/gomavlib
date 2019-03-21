package gomavlib

import (
	"reflect"
	"time"
)

type nodeHeartbeat struct {
	n           *Node
	terminate   chan struct{}
	done        chan struct{}
	heartbeatMp *dialectMessage
	period      time.Duration
}

func newNodeHeartbeat(n *Node, period time.Duration) *nodeHeartbeat {
	// heartbeat message must exist in dialect and correspond to standart heartbeat
	if n.conf.Dialect == nil {
		return nil
	}
	mp, ok := n.conf.Dialect.messages[0]
	if ok == false || mp.crcExtra != 50 {
		return nil
	}

	h := &nodeHeartbeat{
		n:           n,
		terminate:   make(chan struct{}),
		done:        make(chan struct{}),
		heartbeatMp: mp,
		period:      period,
	}
	go h.do()
	return h
}

func (h *nodeHeartbeat) close() {
	h.terminate <- struct{}{}
	<-h.done
}

func (h *nodeHeartbeat) do() {
	defer func() { h.done <- struct{}{} }()

	ticker := time.NewTicker(h.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			msg := reflect.New(h.heartbeatMp.elemType)
			msg.Elem().FieldByName("Type").Set(reflect.ValueOf(uint8(6)))      // MAV_TYPE_GCS
			msg.Elem().FieldByName("Autopilot").Set(reflect.ValueOf(uint8(0))) // MAV_AUTOPILOT_GENERIC
			msg.Elem().FieldByName("BaseMode").Set(reflect.ValueOf(uint8(0)))
			msg.Elem().FieldByName("CustomMode").Set(reflect.ValueOf(uint32(0)))
			msg.Elem().FieldByName("SystemStatus").Set(reflect.ValueOf(uint8(4))) // MAV_STATE_ACTIVE
			msg.Elem().FieldByName("MavlinkVersion").Set(reflect.ValueOf(uint8(3)))
			h.n.WriteMessageAll(msg.Interface().(Message))

		case <-h.terminate:
			return
		}
	}
}
