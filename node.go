// gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go
// programming language. It can power UGVs, UAVs, ground stations, monitoring
// systems or routers acting in a Mavlink network.
//
// Mavlink is a lighweight and transport-independent protocol that is mostly
// used to communicate with unmanned ground vehicles (UGV) and unmanned aerial
// vehicles (UAV, drones, quadcopters, multirotors). It is supported by both
// of the most common open source drone softwares (Ardupilot and PX4).
package gomavlib

import (
	"fmt"
	"io"
	"reflect"
	"sync"
	"time"
)

var signatureReferenceDate = time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)

type Version int

const (
	// V2 means wrap outgoing messages in v2 frames.
	V2 Version = iota
	// V1 means wrap outgoing messages in v1 frames.
	V1
)

type heartbeatTicker struct {
	n           *Node
	terminate   chan struct{}
	done        chan struct{}
	heartbeatMp *messageParser
	period      time.Duration
}

func newHeartbeatTicker(n *Node, period time.Duration) *heartbeatTicker {
	// heartbeat message must exist in dialect and correspond to standart heartbeat
	mp, ok := n.frameParser.messageParsers[0]
	if ok == false || mp.crcExtra != 50 {
		return nil
	}

	h := &heartbeatTicker{
		n:           n,
		terminate:   make(chan struct{}),
		done:        make(chan struct{}),
		heartbeatMp: mp,
		period:      period,
	}
	go h.do()
	return h
}

func (h *heartbeatTicker) close() {
	h.terminate <- struct{}{}
	<-h.done
}

func (h *heartbeatTicker) do() {
	defer func() { h.done <- struct{}{} }()

	ticker := time.NewTicker(h.period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			msg := reflect.New(h.heartbeatMp.elemType)
			msg.Elem().FieldByName("Type").Set(reflect.ValueOf(uint8(6)))      // MAV_TYPE_GCS
			msg.Elem().FieldByName("Autopilot").Set(reflect.ValueOf(uint8(0))) // MAV_AUTOPILOT_GENERIC
			msg.Elem().FieldByName("BaseMode").Set(reflect.ValueOf(uint8(0)))
			msg.Elem().FieldByName("CustomMode").Set(reflect.ValueOf(uint32(0)))
			msg.Elem().FieldByName("SystemStatus").Set(reflect.ValueOf(uint8(4))) // MAV_STATE_ACTIVE
			msg.Elem().FieldByName("MavlinkVersion").Set(reflect.ValueOf(uint8(3)))
			h.n.WriteMessageAll(msg.Interface().(Message))

		case <-h.terminate:
			return
		}
	}
}

// NodeConf allows to configure a Node.
type NodeConf struct {
	// Mavlink version used to encode frames. See Version
	// for the available options.
	Version Version
	// contains the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect []Message
	// these are used to identify this node in the network.
	// They are added to every outgoing message.
	SystemId    byte
	ComponentId byte
	// contains the transport layers with which this node will
	// communicate. Each transport contains one or more channels.
	Transports []TransportConf

	// (optional) the secret key used to verify incoming frames.
	// Non signed frames are discarded. This feature requires Mavlink v2.
	SignatureInKey *FrameSignatureKey
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires Mavlink v2.
	SignatureOutKey *FrameSignatureKey

	// (optional) disable sthe periodic sending of heartbeats to
	// open channels.
	HeartbeatDisable bool
	// (optional) set the period between heartbeats.
	// It defaults to 5 seconds.
	HeartbeatPeriod time.Duration

	// (optional) disables checksum validation of incoming frames.
	// Not recommended, useful only for debugging purposes.
	ChecksumDisable bool
}

// Node represents our node in the network.
type Node struct {
	conf            NodeConf
	mutex           sync.Mutex
	wg              sync.WaitGroup
	chanAccepters   map[transportChannelAccepter]struct{}
	channels        map[*TransportChannel]struct{}
	frameParser     *FrameParser
	frameQueue      chan *NodeReadResult
	writeDone       chan struct{}
	heartbeatTicker *heartbeatTicker
}

// NewNode allocates a Node and connects it to a mavlink network through transports.
// See NodeConf for the options.
func NewNode(conf NodeConf) (*Node, error) {
	if conf.SystemId < 1 {
		return nil, fmt.Errorf("SystemId must be >= 1")
	}
	if conf.ComponentId < 1 {
		return nil, fmt.Errorf("ComponentId must be >= 1")
	}
	if len(conf.Transports) == 0 {
		return nil, fmt.Errorf("at least one transport must be provided")
	}
	if conf.SignatureInKey != nil && conf.Version != V2 {
		return nil, fmt.Errorf("SignatureInKey requires V2 frames")
	}
	if conf.SignatureOutKey != nil && conf.Version != V2 {
		return nil, fmt.Errorf("SignatureOutKey requires V2 frames")
	}
	if conf.HeartbeatPeriod == 0 {
		conf.HeartbeatPeriod = 5 * time.Second
	}

	// init frame parser
	frameParser, err := NewFrameParser(FrameParserConf{
		Dialect: conf.Dialect,
	})
	if err != nil {
		return nil, fmt.Errorf("frame parser: %s", err)
	}

	n := &Node{
		conf:          conf,
		chanAccepters: make(map[transportChannelAccepter]struct{}),
		channels:      make(map[*TransportChannel]struct{}),
		frameParser:   frameParser,
		frameQueue:    make(chan *NodeReadResult),
		writeDone:     make(chan struct{}),
	}

	// init transports
	for _, tconf := range conf.Transports {
		tp, err := tconf.init()
		if err != nil {
			n.Close()
			return nil, err
		}

		if tm, ok := tp.(transportChannelAccepter); ok {
			n.startChannelAccepter(tm)

		} else if ts, ok := tp.(transportChannelSingle); ok {
			n.startChannel(ts)

		} else {
			panic(fmt.Errorf("transport %d does not implement required interface", tp))
		}
	}

	// start heartbeat
	if n.conf.HeartbeatDisable == false {
		n.heartbeatTicker = newHeartbeatTicker(n, n.conf.HeartbeatPeriod)
	}
	return n, nil
}

func (n *Node) Close() {
	func() {
		n.mutex.Lock()
		defer n.mutex.Unlock()

		if n.heartbeatTicker != nil {
			n.heartbeatTicker.close()
		}

		for mc := range n.chanAccepters {
			mc.Close()
		}

		for ch := range n.channels {
			ch.rwc.Close()
		}
	}()

	// consume queued frames up to close(n.frameQueue)
	go func() {
		for {
			_, ok := n.Read()
			if ok == false {
				break
			}
		}
	}()

	n.wg.Wait()

	// close queue after ensuring no one will write to it
	close(n.frameQueue)
}

func (n *Node) startChannelAccepter(tm transportChannelAccepter) {
	n.chanAccepters[tm] = struct{}{}

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		for {
			ch, err := tm.Accept()
			if err != nil {
				if err != errorTerminated {
					panic("errorTerminated is the only error allowed here")
				}
				break
			}

			n.startChannel(ch)
		}
	}()
}

func (n *Node) startChannel(rwc io.ReadWriteCloser) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	conn := &TransportChannel{
		rwc:       rwc,
		writeChan: make(chan interface{}),
	}
	n.channels[conn] = struct{}{}

	// reader
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		defer func() {
			n.mutex.Lock()
			defer n.mutex.Unlock()
			delete(n.channels, conn)
			close(conn.writeChan)
		}()

		buf := make([]byte, netBufferSize)
		for {
			bufn, err := conn.rwc.Read(buf)
			if err != nil {
				// avoid calling twice Close()
				if err != errorTerminated {
					conn.rwc.Close()
				}
				return
			}

			frame, err := n.frameParser.Decode(buf[:bufn], !n.conf.ChecksumDisable, n.conf.SignatureInKey)
			if err != nil {
				fmt.Printf("SKIPPED DUE TO ERR: %v\n", err)
				continue
			}

			n.frameQueue <- &NodeReadResult{frame, conn}
		}
	}()

	// writer
	nextSequenceId := byte(0)
	signatureLinkId := randomByte()
	n.wg.Add(1)
	go func() {
		defer n.wg.Done()

		for {
			what, ok := <-conn.writeChan
			if ok == false {
				return
			}

			var f Frame
			switch wh := what.(type) {
			case Message:
				// SequenceId and SignatureLinkId are unique for each channel
				if n.conf.Version == V1 {
					f = &FrameV1{
						SequenceId:  nextSequenceId,
						SystemId:    n.conf.SystemId,
						ComponentId: n.conf.ComponentId,
						Message:     wh,
					}
				} else {
					f = &FrameV2{
						SequenceId:      nextSequenceId,
						SystemId:        n.conf.SystemId,
						ComponentId:     n.conf.ComponentId,
						Message:         wh,
						SignatureLinkId: signatureLinkId,
						// Timestamp in 10 microsecond units since 1st January 2015 GMT time
						SignatureTimestamp: (uint64(time.Since(signatureReferenceDate)) / 10000),
					}
				}
				nextSequenceId++

			case Frame:
				f = wh
			}

			byt, err := n.frameParser.Encode(f, !n.conf.ChecksumDisable, n.conf.SignatureOutKey)
			if err == nil {
				conn.rwc.Write(byt)
			}

			n.writeDone <- struct{}{}
		}
	}()
}

// NodeReadResult contains the result of node.Read()
type NodeReadResult struct {
	frame            Frame
	transportChannel *TransportChannel
}

// Frame() returns the Frame containing the message.
func (res *NodeReadResult) Frame() Frame {
	return res.frame
}

// Message() returns the message.
func (res *NodeReadResult) Message() Message {
	return res.frame.GetMessage()
}

// SystemId() returns the sender system id.
func (res *NodeReadResult) SystemId() byte {
	return res.frame.GetSystemId()
}

// ComponentId() returns the sender component id.
func (res *NodeReadResult) ComponentId() byte {
	return res.frame.GetComponentId()
}

// Channel() returns the channel used to send the message.
func (res *NodeReadResult) Channel() *TransportChannel {
	return res.transportChannel
}

// Read reads a single message from available channels.
// NodeReadResult contains all the properties of the received frame and message
// (see NodeReadResult for details).
// bool is true whenever the node is still open.
func (n *Node) Read() (*NodeReadResult, bool) {
	res, ok := <-n.frameQueue
	return res, ok
}

// WriteMessageTo write a message to given channel.
func (n *Node) WriteMessageTo(targetChannel *TransportChannel, message Message) {
	n.writeTo(targetChannel, message)
}

// WriteMessageAll write a message to all channels.
func (n *Node) WriteMessageAll(message Message) {
	n.writeAll(message)
}

// WriteMessageExcept write a message to all channels except specified channel.
func (n *Node) WriteMessageExcept(exceptChannel *TransportChannel, message Message) {
	n.writeExcept(exceptChannel, message)
}

// WriteFrameTo write a frame to given channel.
// This function is intended for routing frames to other nodes, since all
// the fields must be filled manually.
func (n *Node) WriteFrameTo(targetChannel *TransportChannel, frame Frame) {
	n.writeTo(targetChannel, frame)
}

// WriteFrameAll write a frame to all channels.
// This function is intended for routing frames to other nodes, since all
// the fields must be filled manually.
func (n *Node) WriteFrameAll(frame Frame) {
	n.writeAll(frame)
}

// WriteFrameExcept write a frame to all channels except specified channel.
// This function is intended for routing frames to other nodes, since all
// the fields must be filled manually.
func (n *Node) WriteFrameExcept(exceptChannel *TransportChannel, frame Frame) {
	n.writeExcept(exceptChannel, frame)
}

func (n *Node) writeTo(targetChannel *TransportChannel, what interface{}) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if _, ok := n.channels[targetChannel]; ok == false {
		return
	}

	// route to channels
	// wait for responses (otherwise transports can be removed before writing)
	targetChannel.writeChan <- what
	<-n.writeDone
}

func (n *Node) writeAll(what interface{}) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	for conn := range n.channels {
		conn.writeChan <- what
	}
	for range n.channels {
		<-n.writeDone
	}
}

func (n *Node) writeExcept(exceptChannel *TransportChannel, what interface{}) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	count := 0
	for conn := range n.channels {
		if conn != exceptChannel {
			count++
			conn.writeChan <- what
		}
	}
	for i := 0; i < count; i++ {
		<-n.writeDone
	}
}
