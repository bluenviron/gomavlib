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
}

func newNodeHeartbeat(n *Node) *nodeHeartbeat {
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
		terminate:   make(chan struct{}, 1),
		done:        make(chan struct{}),
		heartbeatMp: mp,
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

	ticker := time.NewTicker(h.n.conf.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			msg := reflect.New(h.heartbeatMp.elemType)
			msg.Elem().FieldByName("Type").SetInt(6)      // MAV_TYPE_GCS
			msg.Elem().FieldByName("Autopilot").SetInt(0) // MAV_AUTOPILOT_GENERIC
			msg.Elem().FieldByName("BaseMode").SetInt(0)
			msg.Elem().FieldByName("CustomMode").SetUint(0)
			msg.Elem().FieldByName("SystemStatus").SetInt(4) // MAV_STATE_ACTIVE
			msg.Elem().FieldByName("MavlinkVersion").SetUint(3)
			h.n.WriteMessageAll(msg.Interface().(Message))

		case <-h.terminate:
			return
		}
	}
}
