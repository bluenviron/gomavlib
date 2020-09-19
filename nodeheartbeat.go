package gomavlib

import (
	"reflect"
	"time"
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
	mp, ok := n.conf.Dialect.messages[0]
	if ok == false || mp.crcExtra != 50 {
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
			msg := reflect.New(h.n.conf.Dialect.messages[0].elemType)
			msg.Elem().FieldByName("Type").SetInt(int64(h.n.conf.HeartbeatSystemType))
			msg.Elem().FieldByName("Autopilot").SetInt(int64(h.n.conf.HeartbeatAutopilotType))
			msg.Elem().FieldByName("BaseMode").SetInt(0)
			msg.Elem().FieldByName("CustomMode").SetUint(0)
			msg.Elem().FieldByName("SystemStatus").SetInt(4) // MAV_STATE_ACTIVE
			msg.Elem().FieldByName("MavlinkVersion").SetUint(mavlinkVersion)
			h.n.WriteMessageAll(msg.Interface().(Message))

		case <-h.terminate:
			return
		}
	}
}
