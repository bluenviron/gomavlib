/*
Package gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go
programming language. It can power UGVs, UAVs, ground stations, monitoring
systems or routers acting in a Mavlink network.

Mavlink is a lighweight and transport-independent protocol that is mostly used
to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles
(UAV, drones, quadcopters, multirotors). It is supported by the most common
open-source flight controllers (Ardupilot and PX4).

Basic example (more are available at https://github.com/aler9/gomavlib/tree/master/examples):

  package main

  import (
  	"fmt"
  	"github.com/aler9/gomavlib"
  	"github.com/aler9/gomavlib/dialects/ardupilotmega"
  )

  func main() {
  	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
  		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2,
  		OutSystemId: 10,
  	})
  	if err != nil {
  		panic(err)
  	}
  	defer node.Close()

  	for evt := range node.Events() {
  		if frm,ok := evt.(*gomavlib.EventFrame); ok {
  			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())
  		}
  	}
  }

*/
package gomavlib

import (
	"fmt"
	"sync"
	"time"
)

const (
	// constant for ip-based endpoints
	_NET_BUFFER_SIZE      = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
	_NET_CONNECT_TIMEOUT  = 10 * time.Second
	_NET_RECONNECT_PERIOD = 2 * time.Second
	_NET_READ_TIMEOUT     = 60 * time.Second
	_NET_WRITE_TIMEOUT    = 10 * time.Second
)

type goroutinePool sync.WaitGroup

type goroutinePoolRunnable interface {
	run()
}

func (wg *goroutinePool) Start(what goroutinePoolRunnable) {
	(*sync.WaitGroup)(wg).Add(1)
	go func() {
		defer (*sync.WaitGroup)(wg).Done()
		what.run()
	}()
}

func (wg *goroutinePool) Wait() {
	(*sync.WaitGroup)(wg).Wait()
}

// NodeConf allows to configure a Node.
type NodeConf struct {
	// the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect *Dialect

	// (optional) the secret key used to validate incoming frames.
	// Non signed frames are discarded, as well as frames with a version < 2.0.
	InKey *Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemId byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentId byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires a version >= 2.0.
	OutKey *Key

	// (optional) disables the periodic sending of heartbeats to open channels.
	HeartbeatDisable bool
	// (optional) the period between heartbeats. It defaults to 5 seconds.
	HeartbeatPeriod time.Duration
	// (optional) the system type advertised by heartbeats.
	// It defaults to MAV_TYPE_GCS
	HeartbeatSystemType int
	// (optional) the autopilot type advertised by heartbeats.
	// It defaults to MAV_AUTOPILOT_GENERIC
	HeartbeatAutopilotType int

	// (optional) automatically request streams to detected Ardupilot devices,
	// that need an explicit request in order to emit telemetry stream.
	StreamRequestEnable bool
	// (optional) the requested stream frequency in Hz. It defaults to 4.
	StreamRequestFrequency int
}

// Node is a high-level Mavlink encoder and decoder that works with endpoints.
type Node struct {
	conf              NodeConf
	eventsOut         chan Event
	eventsIn          chan eventIn
	pool              goroutinePool
	channelAccepters  map[*channelAccepter]struct{}
	channels          map[*Channel]struct{}
	nodeHeartbeat     *nodeHeartbeat
	nodeStreamRequest *nodeStreamRequest
}

// NewNode allocates a Node. See NodeConf for the options.
func NewNode(conf NodeConf) (*Node, error) {
	if len(conf.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be provided")
	}
	if conf.HeartbeatPeriod == 0 {
		conf.HeartbeatPeriod = 5 * time.Second
	}
	if conf.HeartbeatSystemType == 0 {
		conf.HeartbeatSystemType = 6 // MAV_TYPE_GCS
	}
	if conf.HeartbeatAutopilotType == 0 {
		conf.HeartbeatAutopilotType = 0 // MAV_AUTOPILOT_GENERIC
	}
	if conf.StreamRequestFrequency == 0 {
		conf.StreamRequestFrequency = 4
	}

	// check Parser configuration here, since Parser is created dynamically
	if conf.OutVersion == 0 {
		return nil, fmt.Errorf("OutVersion not provided")
	}
	if conf.OutSystemId < 1 {
		return nil, fmt.Errorf("SystemId must be >= 1")
	}
	if conf.OutComponentId < 1 {
		conf.OutComponentId = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	n := &Node{
		conf: conf,
		// these can be unbuffered as long as eventsIn's goroutine
		// does not write to eventsOut
		eventsOut:        make(chan Event),
		eventsIn:         make(chan eventIn),
		channelAccepters: make(map[*channelAccepter]struct{}),
		channels:         make(map[*Channel]struct{}),
	}

	closeExisting := func() {
		for ca := range n.channels {
			ca.rwc.Close()
		}
		for ca := range n.channelAccepters {
			ca.eca.Close()
		}
	}

	// endpoints
	for _, tconf := range conf.Endpoints {
		tp, err := tconf.init()
		if err != nil {
			closeExisting()
			return nil, err
		}

		switch ttp := tp.(type) {
		case endpointChannelAccepter:
			ca, err := newChannelAccepter(n, ttp)
			if err != nil {
				closeExisting()
				return nil, err
			}

			n.channelAccepters[ca] = struct{}{}

		case endpointChannelSingle:
			ch, err := newChannel(n, ttp, ttp.Label(), ttp)
			if err != nil {
				closeExisting()
				return nil, err
			}

			n.channels[ch] = struct{}{}

		default:
			panic(fmt.Errorf("endpoint %T does not implement any interface", tp))
		}
	}

	// modules
	n.nodeHeartbeat = newNodeHeartbeat(n)
	n.nodeStreamRequest = newNodeStreamRequest(n)

	if n.nodeHeartbeat != nil {
		n.pool.Start(n.nodeHeartbeat)
	}

	if n.nodeStreamRequest != nil {
		n.pool.Start(n.nodeStreamRequest)
	}

	for ch := range n.channels {
		n.pool.Start(ch)
	}

	for ca := range n.channelAccepters {
		n.pool.Start(ca)
	}

	n.pool.Start(n)

	return n, nil
}

func (n *Node) run() {
outer:
	for rawEvt := range n.eventsIn {
		switch evt := rawEvt.(type) {
		case *eventInChannelNew:
			n.channels[evt.ch] = struct{}{}
			n.pool.Start(evt.ch)

		case *eventInChannelClosed:
			delete(n.channels, evt.ch)
			evt.ch.close()

		case *eventInWriteTo:
			if _, ok := n.channels[evt.ch]; ok == false {
				return
			}
			evt.ch.writeChan <- evt.what

		case *eventInWriteAll:
			for ch := range n.channels {
				ch.writeChan <- evt.what
			}

		case *eventInWriteExcept:
			for ch := range n.channels {
				if ch != evt.except {
					ch.writeChan <- evt.what
				}
			}

		case *eventInClose:
			break outer
		}
	}

	// consume events up to close()
	go func() {
		for range n.eventsIn {
		}
	}()

	if n.nodeHeartbeat != nil {
		n.nodeHeartbeat.close()
	}

	if n.nodeStreamRequest != nil {
		n.nodeStreamRequest.close()
	}

	for ca := range n.channelAccepters {
		ca.close()
	}

	for ch := range n.channels {
		ch.close()
	}
}

// Close halts node operations and waits for all routines to return.
func (n *Node) Close() {
	// consume events up to close()
	// in case user is not calling Events()
	go func() {
		for range n.eventsOut {
		}
	}()

	n.eventsIn <- &eventInClose{}
	n.pool.Wait()
	close(n.eventsIn)
	close(n.eventsOut)
}

// Events returns a channel from which receiving events. Possible events are:
//   *EventChannelOpen
//   *EventChannelClose
//   *EventFrame
//   *EventParseError
//   *EventStreamRequested
// See individual events for meaning and content.
func (n *Node) Events() chan Event {
	return n.eventsOut
}

// WriteMessageTo writes a message to given channel.
func (n *Node) WriteMessageTo(channel *Channel, message Message) {
	n.eventsIn <- &eventInWriteTo{channel, message}
}

// WriteMessageAll writes a message to all channels.
func (n *Node) WriteMessageAll(message Message) {
	n.eventsIn <- &eventInWriteAll{message}
}

// WriteMessageExcept writes a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *Channel, message Message) {
	n.eventsIn <- &eventInWriteExcept{exceptChannel, message}
}

// WriteFrameTo writes a frame to given channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameTo(channel *Channel, frame Frame) {
	n.eventsIn <- &eventInWriteTo{channel, frame}
}

// WriteFrameAll writes a frame to all channels.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameAll(frame Frame) {
	n.eventsIn <- &eventInWriteAll{frame}
}

// WriteFrameExcept writes a frame to all channels except specified channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *Channel, frame Frame) {
	n.eventsIn <- &eventInWriteExcept{exceptChannel, frame}
}
