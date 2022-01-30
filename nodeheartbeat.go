package gomavlib

import (
	"reflect"
	"time"

	"github.com/aler9/gomavlib/pkg/msg"
)

type nodeHeartbeat struct {
	n            *Node
	msgHeartbeat msg.Message

	// in
	terminate chan struct{}

	// out
	done chan struct{}
}

func newNodeHeartbeat(n *Node) *nodeHeartbeat {
	// module is disabled
	if n.conf.HeartbeatDisable {
		return nil
	}

	// dialect must be enabled
	if n.conf.Dialect == nil {
		return nil
	}

	// heartbeat message must exist in dialect and correspond to standard
	msgHeartbeat := func() msg.Message {
		for _, m := range n.conf.Dialect.Messages {
			if m.GetID() == 0 {
				return m
			}
		}
		return nil
	}()
	if msgHeartbeat == nil {
		return nil
	}
	mde, err := msg.NewDecEncoder(msgHeartbeat)
	if err != nil || mde.CRCExtra() != 50 {
		return nil
	}

	h := &nodeHeartbeat{
		n:            n,
		msgHeartbeat: msgHeartbeat,
		terminate:    make(chan struct{}),
		done:         make(chan struct{}),
	}

	return h
}

func (h *nodeHeartbeat) close() {
	close(h.terminate)
	<-h.done
}

func (h *nodeHeartbeat) run() {
	defer close(h.done)

	ticker := time.NewTicker(h.n.conf.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m := reflect.New(reflect.TypeOf(h.msgHeartbeat).Elem())
			m.Elem().FieldByName("Type").SetUint(uint64(h.n.conf.HeartbeatSystemType))
			m.Elem().FieldByName("Autopilot").SetUint(uint64(h.n.conf.HeartbeatAutopilotType))
			m.Elem().FieldByName("BaseMode").SetUint(0)
			m.Elem().FieldByName("CustomMode").SetUint(0)
			m.Elem().FieldByName("SystemStatus").SetUint(4) // MAV_STATE_ACTIVE
			m.Elem().FieldByName("MavlinkVersion").SetUint(uint64(h.n.conf.Dialect.Version))
			h.n.WriteMessageAll(m.Interface().(msg.Message))

		case <-h.terminate:
			return
		}
	}
}
