package gomavlib

import (
	"reflect"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

const (
	heartbeatID  = 0
	heartbeatCRC = 50
)

type nodeHeartbeat struct {
	n            *Node
	msgHeartbeat message.Message

	// in
	terminate chan struct{}

	// out
	done chan struct{}
}

func newNodeHeartbeat(n *Node) *nodeHeartbeat {
	// module is disabled
	if n.HeartbeatDisable {
		return nil
	}

	// dialect must be enabled
	if n.Dialect == nil {
		return nil
	}

	// heartbeat message must exist in dialect and correspond to standard
	msgHeartbeat := func() message.Message {
		for _, m := range n.Dialect.Messages {
			if m.GetID() == heartbeatID {
				return m
			}
		}
		return nil
	}()
	if msgHeartbeat == nil {
		return nil
	}

	mde := &message.ReadWriter{
		Message: msgHeartbeat,
	}
	err := mde.Initialize()
	if err != nil || mde.CRCExtra() != heartbeatCRC {
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

	ticker := time.NewTicker(h.n.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m := reflect.New(reflect.TypeOf(h.msgHeartbeat).Elem())
			m.Elem().FieldByName("Type").SetUint(uint64(h.n.HeartbeatSystemType))
			m.Elem().FieldByName("Autopilot").SetUint(uint64(h.n.HeartbeatAutopilotType))
			m.Elem().FieldByName("BaseMode").SetUint(0)
			m.Elem().FieldByName("CustomMode").SetUint(0)
			m.Elem().FieldByName("SystemStatus").SetUint(4) // MAV_STATE_ACTIVE
			m.Elem().FieldByName("MavlinkVersion").SetUint(uint64(h.n.Dialect.Version))
			h.n.WriteMessageAll(m.Interface().(message.Message)) //nolint:errcheck

		case <-h.terminate:
			return
		}
	}
}
