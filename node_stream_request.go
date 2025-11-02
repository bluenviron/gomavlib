package gomavlib

import (
	"reflect"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

const (
	streamRequestPeriod  = 30 * time.Second
	requestDataStreamID  = 66
	requestDataStreamCRC = 148
)

func findMsgRequestDataStream(messages []message.Message) message.Message {
	for _, m := range messages {
		if m.GetID() == requestDataStreamID {
			rw := &message.ReadWriter{Message: m}
			err := rw.Initialize()
			if err != nil || rw.CRCExtra() != requestDataStreamCRC {
				return nil
			}
			return m
		}
	}
	return nil
}

type streamNode struct {
	Channel     *Channel
	SystemID    byte
	ComponentID byte
}

type nodeStreamRequest struct {
	node *Node

	msgHeartbeat         message.Message
	msgRequestDataStream message.Message
	lastRequestsMutex    sync.Mutex
	lastRequests         map[streamNode]time.Time

	// in
	terminate chan struct{}

	// out
	done chan struct{}
}

func (sr *nodeStreamRequest) initialize() error {
	// module is disabled
	if !sr.node.StreamRequestEnable {
		return errSkip
	}

	// dialect must be enabled
	if sr.node.Dialect == nil {
		return errSkip
	}

	sr.msgHeartbeat = findMsgHeartbeat(sr.node.Dialect.Messages)
	if sr.msgHeartbeat == nil {
		return errSkip
	}

	sr.msgRequestDataStream = findMsgRequestDataStream(sr.node.Dialect.Messages)
	if sr.msgRequestDataStream == nil {
		return errSkip
	}

	sr.lastRequests = make(map[streamNode]time.Time)
	sr.terminate = make(chan struct{})
	sr.done = make(chan struct{})

	return nil
}

func (sr *nodeStreamRequest) close() {
	close(sr.terminate)
	<-sr.done
}

func (sr *nodeStreamRequest) run() {
	defer close(sr.done)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		// periodic cleanup
		case now := <-ticker.C:
			func() {
				sr.lastRequestsMutex.Lock()
				defer sr.lastRequestsMutex.Unlock()

				for rnode, t := range sr.lastRequests {
					if now.Sub(t) >= streamRequestPeriod {
						delete(sr.lastRequests, rnode)
					}
				}
			}()

		case <-sr.terminate:
			return
		}
	}
}

func (sr *nodeStreamRequest) onEventFrame(evt *EventFrame) {
	// message must be heartbeat and sender must be an ardupilot device
	if evt.Message().GetID() != 0 ||
		reflect.ValueOf(evt.Message()).Elem().FieldByName("Autopilot").Uint() != 3 {
		return
	}

	rnode := streamNode{
		Channel:     evt.Channel,
		SystemID:    evt.SystemID(),
		ComponentID: evt.ComponentID(),
	}

	// request streams if sender is new or a request has not been sent in some time
	request := false
	func() {
		sr.lastRequestsMutex.Lock()
		defer sr.lastRequestsMutex.Unlock()

		now := time.Now()

		if _, ok := sr.lastRequests[rnode]; !ok {
			sr.lastRequests[rnode] = time.Now()
			request = true
		} else if now.Sub(sr.lastRequests[rnode]) >= streamRequestPeriod {
			request = true
			sr.lastRequests[rnode] = now
		}
	}()

	if request {
		// https://github.com/mavlink/qgroundcontrol/blob/08f400355a8f3acf1dd8ed91f7f1c757323ac182/src
		// /FirmwarePlugin/APM/APMFirmwarePlugin.cc#L626
		streams := []int{
			1,  // common.MAV_DATA_STREAM_RAW_SENSORS,
			2,  // common.MAV_DATA_STREAM_EXTENDED_STATUS,
			3,  // common.MAV_DATA_STREAM_RC_CHANNELS,
			6,  // common.MAV_DATA_STREAM_POSITION,
			10, // common.MAV_DATA_STREAM_EXTRA1,
			11, // common.MAV_DATA_STREAM_EXTRA2,
			12, // common.MAV_DATA_STREAM_EXTRA3,
		}

		for _, stream := range streams {
			m := reflect.New(reflect.TypeOf(sr.msgRequestDataStream).Elem())
			m.Elem().FieldByName("TargetSystem").SetUint(uint64(evt.SystemID()))
			m.Elem().FieldByName("TargetComponent").SetUint(uint64(evt.ComponentID()))
			m.Elem().FieldByName("ReqStreamId").SetUint(uint64(stream))
			m.Elem().FieldByName("ReqMessageRate").SetUint(uint64(sr.node.StreamRequestFrequency))
			m.Elem().FieldByName("StartStop").SetUint(uint64(1))
			sr.node.WriteMessageTo(evt.Channel, m.Interface().(message.Message)) //nolint:errcheck
		}

		sr.node.pushEvent(&EventStreamRequested{
			Channel:     evt.Channel,
			SystemID:    evt.SystemID(),
			ComponentID: evt.ComponentID(),
		})
	}
}
