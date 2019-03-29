/*
Package gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go
programming language. It can power UGVs, UAVs, ground stations, monitoring
systems or routers acting in a Mavlink network.

Mavlink is a lighweight and endpoint-independent protocol that is mostly
used to communicate with unmanned ground vehicles (UGV) and unmanned aerial
vehicles (UAV, drones, quadcopters, multirotors). It is supported by both
of the most common open-source flight controllers (Ardupilot and PX4).

Basic example (more are available at https://github.com/gswly/gomavlib/tree/master/example)

  package main

  import (
  	"fmt"
  	"github.com/gswly/gomavlib"
  	"github.com/gswly/gomavlib/dialects/ardupilotmega"
  )

  func main() {
  	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
  		Dialect:     ardupilotmega.Dialect,
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
	"io"
	"sync"
	"time"
)

const (
	// constant for ip-based endpoints
	netBufferSize      = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
	netConnectTimeout  = 10 * time.Second
	netReconnectPeriod = 2 * time.Second
	netReadTimeout     = 60 * time.Second
	netWriteTimeout    = 10 * time.Second
)

// Version allows to set the frame version used to wrap outgoing messages.
type Version int

const (
	// V2 wrap outgoing messages in v2 frames.
	V2 Version = iota
	// V1 wrap outgoing messages in v1 frames.
	V1
)

// NodeConf allows to configure a Node.
type NodeConf struct {
	// contains the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// contains the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect *Dialect

	// (optional) the secret key used to validate incoming frames.
	// Non signed frames are discarded, as well as frames with a version < v2.
	InSignatureKey *FrameSignatureKey

	// Mavlink version used to encode frames. See Version
	// for the available options.
	OutVersion Version
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemId byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentId byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires Mavlink v2.
	OutSignatureKey *FrameSignatureKey

	// (optional) disables the periodic sending of heartbeats to
	// open channels.
	HeartbeatDisable bool
	// (optional) set the period between heartbeats.
	// It defaults to 5 seconds.
	HeartbeatPeriod time.Duration
}

// Node is a high-level Mavlink encoder and decoder that works with endpoints.
// See NodeConf for the options.
type Node struct {
	conf          NodeConf
	wg            sync.WaitGroup
	chanAccepters map[endpointChannelAccepter]struct{}
	writeDone     chan struct{}
	eventChan     chan Event
	channelsMutex sync.Mutex
	channels      map[*Channel]struct{}
	nodeHeartbeat *nodeHeartbeat
}

// NewNode allocates a Node. See NodeConf for the options.
func NewNode(conf NodeConf) (*Node, error) {
	if conf.OutSystemId < 1 {
		return nil, fmt.Errorf("SystemId must be >= 1")
	}
	if conf.OutComponentId < 1 {
		conf.OutComponentId = 1
	}
	if len(conf.Endpoints) == 0 {
		return nil, fmt.Errorf("at least one endpoint must be provided")
	}
	if conf.OutSignatureKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutSignatureKey requires V2 frames")
	}
	if conf.HeartbeatPeriod == 0 {
		conf.HeartbeatPeriod = 5 * time.Second
	}

	n := &Node{
		conf:          conf,
		chanAccepters: make(map[endpointChannelAccepter]struct{}),
		writeDone:     make(chan struct{}),
		eventChan:     make(chan Event),
		channels:      make(map[*Channel]struct{}),
	}

	for _, tconf := range conf.Endpoints {
		tp, err := tconf.init()
		if err != nil {
			n.Close()
			return nil, err
		}

		if tm, ok := tp.(endpointChannelAccepter); ok {
			n.startChannelAccepter(tm)

		} else if ts, ok := tp.(endpointChannelSingle); ok {
			n.createChannel(ts, ts.Label(), ts)

		} else {
			panic(fmt.Errorf("endpoint %T does not implement any interface", tp))
		}
	}

	if n.conf.HeartbeatDisable == false {
		n.nodeHeartbeat = newNodeHeartbeat(n)
	}
	return n, nil
}

// Close stops node operations and wait for all routines to return.
func (n *Node) Close() {
	if n.nodeHeartbeat != nil {
		n.nodeHeartbeat.close()
	}

	for mc := range n.chanAccepters {
		mc.Close()
	}

	func() {
		n.channelsMutex.Lock()
		defer n.channelsMutex.Unlock()

		for ch := range n.channels {
			ch.rwc.Close()
		}
	}()

	// consume events (in case user is not calling Events()) such that we can
	// call close(n.eventChan)
	go func() {
		for range n.Events() {
		}
	}()

	n.wg.Wait()

	// close queue after ensuring no one will write to it
	close(n.eventChan)
}

func (n *Node) startChannelAccepter(tm endpointChannelAccepter) {
	n.chanAccepters[tm] = struct{}{}

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		for {
			label, rwc, err := tm.Accept()
			if err != nil {
				if err != errorTerminated {
					panic("errorTerminated is the only error allowed here")
				}
				break
			}

			n.createChannel(tm, label, rwc)
		}
	}()
}

func (n *Node) createChannel(e Endpoint, label string, rwc io.ReadWriteCloser) {
	n.channelsMutex.Lock()
	defer n.channelsMutex.Unlock()

	channel := &Channel{
		Endpoint:  e,
		label:     label,
		rwc:       rwc,
		writeChan: make(chan interface{}),
	}
	n.channels[channel] = struct{}{}

	parser, _ := NewParser(ParserConf{
		Reader:             channel.rwc,
		Writer:             channel.rwc,
		Dialect:            n.conf.Dialect,
		InSignatureKey:     n.conf.InSignatureKey,
		OutSystemId:        n.conf.OutSystemId,
		OutComponentId:     n.conf.OutComponentId,
		OutSignatureLinkId: randomByte(),
		OutSignatureKey:    n.conf.OutSignatureKey,
	})

	// reader
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer func() {
			n.channelsMutex.Lock()
			delete(n.channels, channel)
			n.channelsMutex.Unlock()
			close(channel.writeChan)
			n.eventChan <- &EventChannelClose{channel}
		}()

		n.eventChan <- &EventChannelOpen{channel}

		for {
			frame, err := parser.Read()
			if err != nil {
				// continue in case of parse errors
				if _, ok := err.(*ParserError); ok {
					n.eventChan <- &EventParseError{err, channel}
					continue
				}
				// avoid calling twice Close()
				if err != errorTerminated {
					channel.rwc.Close()
				}
				return
			}

			n.eventChan <- &EventFrame{frame, channel}
		}
	}()

	// writer
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		for what := range channel.writeChan {
			switch wh := what.(type) {
			case Message:
				if n.conf.OutVersion == V1 {
					parser.Write(&FrameV1{Message: wh}, false)
				} else {
					parser.Write(&FrameV2{Message: wh}, false)
				}

			case Frame:
				parser.Write(wh, true)
			}

			n.writeDone <- struct{}{}
		}
	}()
}

// Events returns a channel from which receiving events. Possible events are:
//   *EventChannelOpen
//   *EventChannelClose
//   *EventFrame
//   *EventParseError
// See individual events for meaning and content.
func (n *Node) Events() chan Event {
	return n.eventChan
}

// WriteMessageTo write a message to given channel.
func (n *Node) WriteMessageTo(channel *Channel, message Message) {
	n.writeTo(channel, message)
}

// WriteMessageAll write a message to all channels.
func (n *Node) WriteMessageAll(message Message) {
	n.writeAll(message)
}

// WriteMessageExcept write a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *Channel, message Message) {
	n.writeExcept(exceptChannel, message)
}

// WriteFrameTo write a frame to given channel.
// This function is intended for routing frames to other nodes, since all
// fields must be filled manually.
func (n *Node) WriteFrameTo(channel *Channel, frame Frame) {
	n.writeTo(channel, frame)
}

// WriteFrameAll write a frame to all channels.
// This function is intended for routing frames to other nodes, since all
// fields must be filled manually.
func (n *Node) WriteFrameAll(frame Frame) {
	n.writeAll(frame)
}

// WriteFrameExcept write a frame to all channels except specified channel.
// This function is intended for routing frames to other nodes, since all
// fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *Channel, frame Frame) {
	n.writeExcept(exceptChannel, frame)
}

func (n *Node) writeTo(channel *Channel, what interface{}) {
	n.channelsMutex.Lock()
	defer n.channelsMutex.Unlock()

	if _, ok := n.channels[channel]; ok == false {
		return
	}

	// route to channels
	// wait for responses (otherwise endpoints can be removed before writing)
	channel.writeChan <- what
	<-n.writeDone
}

func (n *Node) writeAll(what interface{}) {
	n.channelsMutex.Lock()
	defer n.channelsMutex.Unlock()

	for channel := range n.channels {
		channel.writeChan <- what
		defer func() { <-n.writeDone }()
	}
}

func (n *Node) writeExcept(exceptChannel *Channel, what interface{}) {
	n.channelsMutex.Lock()
	defer n.channelsMutex.Unlock()

	for channel := range n.channels {
		if channel != exceptChannel {
			channel.writeChan <- what
			defer func() { <-n.writeDone }()
		}
	}
}
