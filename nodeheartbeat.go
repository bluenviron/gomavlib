package gomavlib

import (
	"reflect"
	"time"

	"github.com/aler9/gomavlib/msg"
)

type nodeHeartbeat struct {
	n *Node

	terminate chan struct{}
	done      chan struct{}
}

func newNodeHeartbeat(n *Node) *nodeHeartbeat {
	// module is disabled
	if n.conf.HeartbeatDisable == true {
		return nil
	}

	// dialect must be enabled
	if n.conf.Dialect == nil {
		return nil
	}

	// heartbeat message must exist in dialect and correspond to standard
	mp, ok := n.conf.Dialect.messageDEs[0]
	if ok == false || mp.CRCExtra != 50 {
		return nil
	}

	h := &nodeHeartbeat{
		n:         n,
		terminate: make(chan struct{}),
		done:      make(chan struct{}),
	}

	return h
}

func (h *nodeHeartbeat) close() {
	close(h.terminate)
	<-h.done
}

func (h *nodeHeartbeat) run() {
	defer close(h.done)

	// take version from dialect if possible
	mavlinkVersion := uint64(3)
	if h.n.conf.Dialect != nil {
		mavlinkVersion = uint64(h.n.conf.Dialect.version)
	}

	ticker := time.NewTicker(h.n.conf.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m := reflect.New(h.n.conf.Dialect.messageDEs[0].ElemType)
			m.Elem().FieldByName("Type").SetInt(int64(h.n.conf.HeartbeatSystemType))
			m.Elem().FieldByName("Autopilot").SetInt(int64(h.n.conf.HeartbeatAutopilotType))
			m.Elem().FieldByName("BaseMode").SetInt(0)
			m.Elem().FieldByName("CustomMode").SetUint(0)
			m.Elem().FieldByName("SystemStatus").SetInt(4) // MAV_STATE_ACTIVE
			m.Elem().FieldByName("MavlinkVersion").SetUint(mavlinkVersion)
			h.n.WriteMessageAll(m.Interface().(msg.Message))

		case <-h.terminate:
			return
		}
	}
}
