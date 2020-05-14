package gomavlib

import (
	"reflect"
	"sync"
	"time"
)

const (
	_STREAM_REQUEST_PERIOD = 30 * time.Second
)

type streamNode struct {
	Channel     *Channel
	SystemId    byte
	ComponentId byte
}

type nodeStreamRequest struct {
	n                 *Node
	terminate         chan struct{}
	lastRequestsMutex sync.Mutex
	lastRequests      map[streamNode]time.Time
}

func newNodeStreamRequest(n *Node) *nodeStreamRequest {
	// module is disabled
	if n.conf.StreamRequestEnable == false {
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

	// request data stream message must exist in dialect and correspond to standard
	mp, ok = n.conf.Dialect.messages[66]
	if ok == false || mp.crcExtra != 148 {
		return nil
	}

	sr := &nodeStreamRequest{
		n:            n,
		terminate:    make(chan struct{}),
		lastRequests: make(map[streamNode]time.Time),
	}

	return sr
}

func (sr *nodeStreamRequest) close() {
	close(sr.terminate)
}

func (sr *nodeStreamRequest) run() {
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
					if now.Sub(t) >= _STREAM_REQUEST_PERIOD {
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
	if evt.Message().GetId() != 0 ||
		reflect.ValueOf(evt.Message()).Elem().FieldByName("Autopilot").Int() != 3 {
		return
	}

	rnode := streamNode{
		Channel:     evt.Channel,
		SystemId:    evt.SystemId(),
		ComponentId: evt.ComponentId(),
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

		} else {
			if now.Sub(sr.lastRequests[rnode]) >= _STREAM_REQUEST_PERIOD {
				request = true
				sr.lastRequests[rnode] = now
			}
		}
	}()

	if request == true {
		// https://github.com/mavlink/qgroundcontrol/blob/08f400355a8f3acf1dd8ed91f7f1c757323ac182/src/FirmwarePlugin/APM/APMFirmwarePlugin.cc#L626
		streams := []int{
			1,  //common.MAV_DATA_STREAM_RAW_SENSORS,
			2,  //common.MAV_DATA_STREAM_EXTENDED_STATUS,
			3,  //common.MAV_DATA_STREAM_RC_CHANNELS,
			6,  //common.MAV_DATA_STREAM_POSITION,
			10, //common.MAV_DATA_STREAM_EXTRA1,
			11, //common.MAV_DATA_STREAM_EXTRA2,
			12, //common.MAV_DATA_STREAM_EXTRA3,
		}

		for _, stream := range streams {
			msg := reflect.New(sr.n.conf.Dialect.messages[66].elemType)
			msg.Elem().FieldByName("TargetSystem").SetUint(uint64(evt.SystemId()))
			msg.Elem().FieldByName("TargetComponent").SetUint(uint64(evt.ComponentId()))
			msg.Elem().FieldByName("ReqStreamId").SetUint(uint64(stream))
			msg.Elem().FieldByName("ReqMessageRate").SetUint(uint64(sr.n.conf.StreamRequestFrequency))
			msg.Elem().FieldByName("StartStop").SetUint(uint64(1))
			sr.n.WriteMessageTo(evt.Channel, msg.Interface().(Message))
		}

		sr.n.eventsOut <- &EventStreamRequested{
			Channel:     evt.Channel,
			SystemId:    evt.SystemId(),
			ComponentId: evt.ComponentId(),
		}
	}
}
