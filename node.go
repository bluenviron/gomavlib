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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

var (
	errTerminated = fmt.Errorf("terminated")
	errSkip       = fmt.Errorf("skip")
)

type writeToReq struct {
	ch   *Channel
	what any
}

type writeExceptReq struct {
	except *Channel
	what   any
}

// NodeConf allows to configure a Node.
//
// Deprecated: configuration has been moved inside Node.
type NodeConf struct {
	// endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// (optional) dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect *dialect.Dialect

	// (optional) secret key used to validate incoming frames.
	// Non signed frames are discarded, as well as frames with a version < 2.0.
	InKey *frame.V2Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires a version >= 2.0.
	OutKey *frame.V2Key

	// (optional) disables the periodic sending of heartbeats to open channels.
	HeartbeatDisable bool
	// (optional) period between heartbeats. It defaults to 5 seconds.
	HeartbeatPeriod time.Duration
	// (optional) system type advertised by heartbeats.
	// It defaults to MAV_TYPE_GCS
	HeartbeatSystemType int
	// (optional) autopilot type advertised by heartbeats.
	// It defaults to MAV_AUTOPILOT_GENERIC
	HeartbeatAutopilotType int

	// (optional) automatically request streams to detected Ardupilot devices,
	// that need an explicit request in order to emit telemetry stream.
	StreamRequestEnable bool
	// (optional) requested stream frequency in Hz. It defaults to 4.
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

// NewNode allocates a Node. See NodeConf for the options.
//
// Deprecated: replaced by Node.Initialize().
func NewNode(conf NodeConf) (*Node, error) {
	n := &Node{
		Endpoints:              conf.Endpoints,
		Dialect:                conf.Dialect,
		InKey:                  conf.InKey,
		OutVersion:             conf.OutVersion,
		OutSystemID:            conf.OutSystemID,
		OutComponentID:         conf.OutComponentID,
		OutKey:                 conf.OutKey,
		HeartbeatDisable:       conf.HeartbeatDisable,
		HeartbeatPeriod:        conf.HeartbeatPeriod,
		HeartbeatSystemType:    conf.HeartbeatSystemType,
		HeartbeatAutopilotType: conf.HeartbeatAutopilotType,
		StreamRequestEnable:    conf.StreamRequestEnable,
		StreamRequestFrequency: conf.StreamRequestFrequency,
		ReadTimeout:            conf.ReadTimeout,
		WriteTimeout:           conf.WriteTimeout,
		IdleTimeout:            conf.IdleTimeout,
	}
	err := n.Initialize()
	return n, err
}

// Node is a high-level Mavlink encoder and decoder that works with endpoints.
type Node struct {
	// endpoints with which this node will
	// communicate. Each endpoint contains zero or more channels
	Endpoints []EndpointConf

	// (optional) dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect *dialect.Dialect

	// (optional) secret key used to validate incoming frames.
	// Non signed frames are discarded, as well as frames with a version < 2.0.
	InKey *frame.V2Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires a version >= 2.0.
	OutKey *frame.V2Key

	// (optional) disables the periodic sending of heartbeats to open channels.
	HeartbeatDisable bool
	// (optional) period between heartbeats. It defaults to 5 seconds.
	HeartbeatPeriod time.Duration
	// (optional) system type advertised by heartbeats.
	// It defaults to MAV_TYPE_GCS
	HeartbeatSystemType int
	// (optional) autopilot type advertised by heartbeats.
	// It defaults to MAV_AUTOPILOT_GENERIC
	HeartbeatAutopilotType int

	// (optional) automatically request streams to detected Ardupilot devices,
	// that need an explicit request in order to emit telemetry stream.
	StreamRequestEnable bool
	// (optional) requested stream frequency in Hz. It defaults to 4.
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

	//
	// private
	//

	dialectRW         *dialect.ReadWriter
	wg                sync.WaitGroup
	channelProviders  map[*channelProvider]struct{}
	channels          map[*Channel]struct{}
	nodeHeartbeat     *nodeHeartbeat
	nodeStreamRequest *nodeStreamRequest

	// in
	chNewChannel   chan *Channel
	chCloseChannel chan *Channel
	chWriteTo      chan writeToReq
	chWriteAll     chan any
	chWriteExcept  chan writeExceptReq
	terminate      chan struct{}

	// out
	chEvent chan Event
	done    chan struct{}
}

// Initialize initializes a Node.
func (n *Node) Initialize() error {
	if len(n.Endpoints) == 0 {
		return fmt.Errorf("at least one endpoint must be provided")
	}
	if n.HeartbeatPeriod == 0 {
		n.HeartbeatPeriod = 5 * time.Second
	}
	if n.HeartbeatSystemType == 0 {
		n.HeartbeatSystemType = 6 // MAV_TYPE_GCS
	}
	if n.HeartbeatAutopilotType == 0 {
		n.HeartbeatAutopilotType = 0 // MAV_AUTOPILOT_GENERIC
	}
	if n.StreamRequestFrequency == 0 {
		n.StreamRequestFrequency = 4
	}

	// check Transceiver configuration here, since Transceiver is created dynamically
	if n.OutVersion == 0 {
		return fmt.Errorf("OutVersion not provided")
	}
	if n.OutSystemID < 1 {
		return fmt.Errorf("OutSystemID must be greater than one")
	}
	if n.OutComponentID < 1 {
		n.OutComponentID = 1
	}
	if n.OutKey != nil && n.OutVersion != V2 {
		return fmt.Errorf("OutKey requires V2 frames")
	}

	if n.ReadTimeout == 0 {
		n.ReadTimeout = 10 * time.Second
	}
	if n.WriteTimeout == 0 {
		n.WriteTimeout = 10 * time.Second
	}
	if n.IdleTimeout == 0 {
		n.IdleTimeout = 60 * time.Second
	}

	var dialectRW *dialect.ReadWriter
	if n.Dialect != nil {
		dialectRW = &dialect.ReadWriter{Dialect: n.Dialect}
		err := dialectRW.Initialize()
		if err != nil {
			return err
		}
	}

	n.dialectRW = dialectRW
	n.channelProviders = make(map[*channelProvider]struct{})
	n.channels = make(map[*Channel]struct{})
	n.chNewChannel = make(chan *Channel)
	n.chCloseChannel = make(chan *Channel)
	n.chWriteTo = make(chan writeToReq)
	n.chWriteAll = make(chan any)
	n.chWriteExcept = make(chan writeExceptReq)
	n.terminate = make(chan struct{})
	n.chEvent = make(chan Event)
	n.done = make(chan struct{})

	closeExisting := func() {
		for ca := range n.channelProviders {
			ca.close()
		}
	}

	// endpoints
	for _, conf := range n.Endpoints {
		endpoint, err := conf.init(n)
		if err != nil {
			closeExisting()
			return err
		}

		ca := &channelProvider{
			node:     n,
			endpoint: endpoint,
		}
		err = ca.initialize()
		if err != nil {
			closeExisting()
			return err
		}

		n.channelProviders[ca] = struct{}{}
	}

	n.nodeHeartbeat = &nodeHeartbeat{
		node: n,
	}
	err := n.nodeHeartbeat.initialize()
	if err != nil {
		if errors.Is(err, errSkip) {
			n.nodeHeartbeat = nil
		} else {
			return err
		}
	}

	n.nodeStreamRequest = &nodeStreamRequest{
		node: n,
	}
	err = n.nodeStreamRequest.initialize()
	if err != nil {
		if errors.Is(err, errSkip) {
			n.nodeStreamRequest = nil
		} else {
			return err
		}
	}

	if n.nodeHeartbeat != nil {
		go n.nodeHeartbeat.run()
	}

	if n.nodeStreamRequest != nil {
		go n.nodeStreamRequest.run()
	}

	for ca := range n.channelProviders {
		ca.start()
	}

	go n.run()

	return nil
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
			req.ch.write(req.what)

		case what := <-n.chWriteAll:
			for ch := range n.channels {
				ch.write(what)
			}

		case req := <-n.chWriteExcept:
			for ch := range n.channels {
				if ch != req.except {
					ch.write(req.what)
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

	for ca := range n.channelProviders {
		ca.close()
	}

	for ch := range n.channels {
		ch.close()
	}

	n.wg.Wait()

	close(n.chEvent)
}

// FixFrame recomputes the Frame checksum and signature.
// This can be called on Frames whose content has been edited.
func (n *Node) FixFrame(fr frame.Frame) error {
	err := n.encodeFrame(fr)
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
	if ff, ok := fr.(*frame.V2Frame); ok && n.OutKey != nil {
		ff.Signature = ff.GenerateSignature(n.OutKey)
	}

	return nil
}

func (n *Node) encodeFrame(fr frame.Frame) error {
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if n.dialectRW == nil {
			return fmt.Errorf("dialect is nil")
		}

		mp := n.dialectRW.GetMessage(fr.GetMessage().GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		_, isV2 := fr.(*frame.V2Frame)
		msgRaw := mp.Write(fr.GetMessage(), isV2)

		switch fr := fr.(type) {
		case *frame.V1Frame:
			fr.Message = msgRaw
		case *frame.V2Frame:
			fr.Message = msgRaw
		}
	}

	return nil
}

func (n *Node) encodeMessage(msg message.Message) (message.Message, error) {
	if _, ok := msg.(*message.MessageRaw); !ok {
		if n.dialectRW == nil {
			return nil, fmt.Errorf("dialect is nil")
		}

		mp := n.dialectRW.GetMessage(msg.GetID())
		if mp == nil {
			return nil, fmt.Errorf("message is not in the dialect")
		}

		msgRaw := mp.Write(msg, n.OutVersion == V2)
		return msgRaw, nil
	}

	return msg, nil
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
func (n *Node) WriteMessageTo(channel *Channel, m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteTo <- writeToReq{channel, m}:
	case <-n.terminate:
	}

	return nil
}

// WriteMessageAll writes a message to all channels.
func (n *Node) WriteMessageAll(m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteAll <- m:
	case <-n.terminate:
	}

	return nil
}

// WriteMessageExcept writes a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *Channel, m message.Message) error {
	m, err := n.encodeMessage(m)
	if err != nil {
		return err
	}

	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, m}:
	case <-n.terminate:
	}

	return nil
}

// WriteFrameTo writes a frame to given channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameTo(channel *Channel, fr frame.Frame) error {
	err := n.encodeFrame(fr)
	if err != nil {
		return err
	}

	select {
	case n.chWriteTo <- writeToReq{channel, fr}:
	case <-n.terminate:
	}

	return nil
}

// WriteFrameAll writes a frame to all channels.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameAll(fr frame.Frame) error {
	err := n.encodeFrame(fr)
	if err != nil {
		return err
	}

	select {
	case n.chWriteAll <- fr:
	case <-n.terminate:
	}

	return nil
}

// WriteFrameExcept writes a frame to all channels except specified channel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *Channel, fr frame.Frame) error {
	err := n.encodeFrame(fr)
	if err != nil {
		return err
	}

	select {
	case n.chWriteExcept <- writeExceptReq{exceptChannel, fr}:
	case <-n.terminate:
	}

	return nil
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
