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
	node *Node

	msgHeartbeat message.Message

	// in
	terminate chan struct{}

	// out
	done chan struct{}
}

func (h *nodeHeartbeat) initialize() error {
	// module is disabled
	if h.node.HeartbeatDisable {
		return errSkip
	}

	// dialect must be enabled
	if h.node.Dialect == nil {
		return errSkip
	}

	// heartbeat message must exist in dialect and correspond to standard
	h.msgHeartbeat = func() message.Message {
		for _, m := range h.node.Dialect.Messages {
			if m.GetID() == heartbeatID {
				return m
			}
		}
		return nil
	}()
	if h.msgHeartbeat == nil {
		return errSkip
	}

	mde := &message.ReadWriter{
		Message: h.msgHeartbeat,
	}
	err := mde.Initialize()
	if err != nil || mde.CRCExtra() != heartbeatCRC {
		return errSkip
	}

	h.terminate = make(chan struct{})
	h.done = make(chan struct{})

	return nil
}

func (h *nodeHeartbeat) close() {
	close(h.terminate)
	<-h.done
}

func (h *nodeHeartbeat) run() {
	defer close(h.done)

	ticker := time.NewTicker(h.node.HeartbeatPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m := reflect.New(reflect.TypeOf(h.msgHeartbeat).Elem())
			m.Elem().FieldByName("Type").SetUint(uint64(h.node.HeartbeatSystemType))
			m.Elem().FieldByName("Autopilot").SetUint(uint64(h.node.HeartbeatAutopilotType))
			m.Elem().FieldByName("BaseMode").SetUint(0)
			m.Elem().FieldByName("CustomMode").SetUint(0)
			m.Elem().FieldByName("SystemStatus").SetUint(4) // MAV_STATE_ACTIVE
			m.Elem().FieldByName("MavlinkVersion").SetUint(uint64(h.node.Dialect.Version))
			h.node.WriteMessageAll(m.Interface().(message.Message)) //nolint:errcheck

		case <-h.terminate:
			return
		}
	}
}
