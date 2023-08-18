/*
Package gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go
programming language. It can power UGVs, UAVs, ground stations, monitoring
systems or routers acting in a Mavlink network.

Mavlink is a lighweight and transport-independent protocol that is mostly used
to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles
(UAV, drones, quadcopters, multirotors). It is supported by the most common
open-source flight controllers (Ardupilot and PX4).

Examples are available at https://github.com/bluenviron/gomavlib/tree/main/examples
*/
package gomavlib

import (
	"fmt"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/frame"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

var errTerminated = fmt.Errorf("terminated")

type writeToReq struct {
	ch   *Channel
	what interface{}
}

type writeExceptReq struct {
	except *Channel
	what   interface{}
}

// NodeConf allows to configure a Node.
type NodeConf struct {
	// the endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect *dialect.Dialect

	// (optional) the secret key used to validate incoming frames.
	// Non signed frames are discarded, as well as frames with a version < 2.0.
	InKey *frame.V2Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires a version >= 2.0.
	OutKey *frame.V2Key

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

	// (optional) read timeout.
	// It defaults to 10 seconds.
	ReadTimeout time.Duration
	// (optional) write timeout.
	// It defaults to 10 seconds.
	WriteTimeout time.Duration
	// (optional) timeout before closing idle connections.
	// It defaults to 60 seconds.
	IdleTimeout time.Duration
}

// Node is a high-level Mavlink encoder and decoder that works with endpoints.
type Node struct {
	conf               NodeConf
	dialectRW          *dialect.ReadWriter
	channelAccepters   map[*channelAccepter]struct{}
	channelAcceptersWg sync.WaitGroup
	channels           map[*Channel]struct{}
	channelsWg         sync.WaitGroup
	nodeHeartbeat      *nodeHeartbeat
	nodeStreamRequest  *nodeStreamRequest

	// in
	chNewChannel   chan *Channel
	chCloseChannel chan *Channel
	chWriteTo      chan writeToReq
	chWriteAll     chan interface{}
	chWriteExcept  chan writeExceptReq
	terminate      chan struct{}

	// out
	chEvent chan Event
	done    chan struct{}
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

	// check Transceiver configuration here, since Transceiver is created dynamically
	if conf.OutVersion == 0 {
		return nil, fmt.Errorf("OutVersion not provided")
	}
	if conf.OutSystemID < 1 {
		return nil, fmt.Errorf("OutSystemID must be greater than one")
	}
	if conf.OutComponentID < 1 {
		conf.OutComponentID = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = 10 * time.Second
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = 10 * time.Second
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 60 * time.Second
	}

	dialectRW, err := func() (*dialect.ReadWriter, error) {
		if conf.Dialect == nil {
			return nil, nil
		}
		return dialect.NewReadWriter(conf.Dialect)
	}()
	if err != nil {
		return nil, err
	}

	n := &Node{
		conf:             conf,
		dialectRW:        dialectRW,
		channelAccepters: make(map[*channelAccepter]struct{}),
		channels:         make(map[*Channel]struct{}),
		chNewChannel:     make(chan *Channel),
		chCloseChannel:   make(chan *Channel),
		chWriteTo:        make(chan writeToReq),
		chWriteAll:       make(chan interface{}),
		chWriteExcept:    make(chan writeExceptReq),
		terminate:        make(chan struct{}),
		chEvent:          make(chan Event),
		done:             make(chan struct{}),
	}

	closeExisting := func() {
		for ch := range n.channels {
			ch.close()
		}
		for ca := range n.channelAccepters {
			ca.close()
		}
	}

	// endpoints
	for _, tconf := range conf.Endpoints {
		tp, err := tconf.init(n)
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
			ch, err := newChannel(n, ttp, ttp.label(), ttp)
			if err != nil {
				closeExisting()
				return nil, err
			}

			n.channels[ch] = struct{}{}

		default:
			panic(fmt.Errorf("endpoint %T does not implement any interface", tp))
		}
	}

	n.nodeHeartbeat = newNodeHeartbeat(n)
	n.nodeStreamRequest = newNodeStreamRequest(n)

	if n.nodeHeartbeat != nil {
		go n.nodeHeartbeat.run()
	}

	if n.nodeStreamRequest != nil {
		go n.nodeStreamRequest.run()
	}

	for ch := range n.channels {
		ch.start()
	}

	for ca := range n.channelAccepters {
		ca.start()
	}

	go n.run()

	return n, nil
}

// Close halts node operations and waits for all routines to return.
func (n *Node) Close() {
	close(n.terminate)
	<-n.done
}

func (n *Node) run() {
	defer close(n.done)

outer:
	for {
		select {
		case ch := <-n.chNewChannel:
			n.channels[ch] = struct{}{}
			ch.start()

		case ch := <-n.chCloseChannel:
			delete(n.channels, ch)

		case req := <-n.chWriteTo:
			if _, ok := n.channels[req.ch]; !ok {
				continue
			}

			var err error
			req.what, err = n.encodeMessage(req.what)
			if err == nil {
				req.ch.write(req.what)
			}

		case what := <-n.chWriteAll:
			var err error
			what, err = n.encodeMessage(what)
			if err == nil {
				for ch := range n.channels {
					ch.write(what)
				}
			}

		case req := <-n.chWriteExcept:
			var err error
			req.what, err = n.encodeMessage(req.what)
			if err == nil {
				for ch := range n.channels {
					if ch != req.except {
						ch.write(req.what)
					}
				}
			}

		case <-n.terminate:
			break outer
		}
	}

	if n.nodeHeartbeat != nil {
		n.nodeHeartbeat.close()
	}

	if n.nodeStreamRequest != nil {
		n.nodeStreamRequest.close()
	}

	for ca := range n.channelAccepters {
		ca.close()
	}
	n.channelAcceptersWg.Wait()

	for ch := range n.channels {
		ch.close()
	}
	n.channelsWg.Wait()

	close(n.chEvent)
}

// FixFrame recomputes the Frame checksum and signature.
// This can be called on Frames whose content has been edited.
func (n *Node) FixFrame(fr frame.Frame) error {
	_, err := n.encodeMessage(fr)
	if err != nil {
		return err
	}

	if n.dialectRW == nil {
		return fmt.Errorf("dialect is nil")
	}

	mp := n.dialectRW.GetMessage(fr.GetMessage().GetID())
	if mp == nil {
		return fmt.Errorf("message is not in the dialect")
	}

	// fill checksum
	switch ff := fr.(type) {
	case *frame.V1Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	case *frame.V2Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	}

	// fill Signature if v2
	if ff, ok := fr.(*frame.V2Frame); ok && n.conf.OutKey != nil {
		ff.Signature = ff.GenerateSignature(n.conf.OutKey)
	}

	return nil
}

// encode messages once before sending them to the channel's frame.ReadWriter.
func (n *Node) encodeMessage(what interface{}) (interface{}, error) {
	switch twhat := what.(type) {
	case message.Message:
		if _, ok := twhat.(*message.MessageRaw); !ok {
			if n.dialectRW == nil {
				return nil, fmt.Errorf("dialect is nil")
			}

			mp := n.dialectRW.GetMessage(twhat.GetID())
			if mp == nil {
				return nil, fmt.Errorf("message is not in the dialect")
			}

			msgRaw := mp.Write(twhat, n.conf.OutVersion == V2)

			return msgRaw, nil
		}

	case frame.Frame:
		if _, ok := twhat.GetMessage().(*message.MessageRaw); !ok {
			if n.dialectRW == nil {
				return nil, fmt.Errorf("dialect is nil")
			}

			mp := n.dialectRW.GetMessage(twhat.GetMessage().GetID())
			if mp == nil {
				return nil, fmt.Errorf("message is not in the dialect")
			}

			_, isV2 := twhat.(*frame.V2Frame)
			msgRaw := mp.Write(twhat.GetMessage(), isV2)

			switch ff := twhat.(type) {
			case *frame.V1Frame:
				ff.Message = msgRaw
			case *frame.V2Frame:
				ff.Message = msgRaw
			}
		}
	}

	return what, nil
}

// Events returns a channel from which receiving events. Possible events are:
//
// * EventChannelOpen
// * EventChannelClose
// * EventFrame
// * EventParseError
// * EventStreamRequested
//
// See individual events for details.
func (n *Node) Events() chan Event {
	return n.chEvent
}

// WriteMessageTo writes a message to given channel.
func (n *Node) WriteMessageTo(channel *Channel, m message.Message) {
	select {
	case n.chWriteTo <- writeToReq{channel, m}:
	case <-n.terminate:
	}
}

// WriteMessageAll writes a message to all channels.
func (n *Node) WriteMessageAll(m message.Message) {
	select {
	case n.chWriteAll <- m:
	case <-n.terminate:
	}
}

// WriteMessageExcept writes a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *Channel, m message.Message) {
	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, m}:
	case <-n.terminate:
	}
}

// WriteFrameTo writes a frame to given channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameTo(channel *Channel, fr frame.Frame) {
	select {
	case n.chWriteTo <- writeToReq{channel, fr}:
	case <-n.terminate:
	}
}

// WriteFrameAll writes a frame to all channels.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameAll(fr frame.Frame) {
	select {
	case n.chWriteAll <- fr:
	case <-n.terminate:
	}
}

// WriteFrameExcept writes a frame to all channels except specified channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *Channel, fr frame.Frame) {
	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, fr}:
	case <-n.terminate:
	}
}

func (n *Node) pushEvent(evt Event) {
	select {
	case n.chEvent <- evt:
	case <-n.terminate:
	}
}

func (n *Node) newChannel(ch *Channel) {
	select {
	case n.chNewChannel <- ch:
	case <-n.terminate:
		ch.close()
	}
}

func (n *Node) closeChannel(ch *Channel) {
	select {
	case n.chCloseChannel <- ch:
	case <-n.terminate:
	}
}
