package gomavlib

import (
	"reflect"
	"time"
)

type nodeHeartbeat struct {
	n           *Node
	terminate   chan struct{}
	heartbeatMp *dialectMessage
}

func newNodeHeartbeat(n *Node) *nodeHeartbeat {
	// heartbeat is disabled
	if n.conf.HeartbeatDisable == true {
		return nil
	}

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
		heartbeatMp: mp,
	}

	return h
}

func (h *nodeHeartbeat) start() {
	h.n.wg.Add(1)

	go func() {
		defer h.n.wg.Done()

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
				msg := reflect.New(h.heartbeatMp.elemType)
				msg.Elem().FieldByName("Type").SetInt(int64(h.n.conf.HeartbeatSystemType))
				msg.Elem().FieldByName("Autopilot").SetInt(0) // MAV_AUTOPILOT_GENERIC
				msg.Elem().FieldByName("BaseMode").SetInt(0)
				msg.Elem().FieldByName("CustomMode").SetUint(0)
				msg.Elem().FieldByName("SystemStatus").SetInt(4) // MAV_STATE_ACTIVE
				msg.Elem().FieldByName("MavlinkVersion").SetUint(mavlinkVersion)
				h.n.WriteMessageAll(msg.Interface().(Message))

			case <-h.terminate:
				return
			}
		}
	}()
}

func (h *nodeHeartbeat) close() {
	h.terminate <- struct{}{}
}
