package gomavlib

import (
	"crypto/rand"
	"io"

	"github.com/bluenviron/gomavlib/v2/pkg/frame"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
	"github.com/bluenviron/gomavlib/v2/pkg/ringbuffer"
)

const (
	// this is low in order to avoid accumulating messages
	// when a channel is reconnecting
	writeBufferSize = 8
)

func randomByte() (byte, error) {
	var buf [1]byte
	_, err := rand.Read(buf[:])
	return buf[0], err
}

// Channel is a communication channel created by an Endpoint.
// An Endpoint can create channels.
// For instance, a TCP client endpoint creates a single channel, while a TCP
// server endpoint creates a channel for each incoming connection.
type Channel struct {
	e           Endpoint
	label       string
	rwc         io.Closer
	n           *Node
	frw         *frame.ReadWriter
	running     bool
	writeBuffer *ringbuffer.RingBuffer

	// in
	terminate chan struct{}
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) (*Channel, error) {
	linkID, err := randomByte()
	if err != nil {
		return nil, err
	}

	frw, err := frame.NewReadWriter(frame.ReadWriterConf{
		ReadWriter:  rwc,
		DialectRW:   n.dialectRW,
		InKey:       n.conf.InKey,
		OutSystemID: n.conf.OutSystemID,
		OutVersion: func() frame.WriterOutVersion {
			if n.conf.OutVersion == V2 {
				return frame.V2
			}
			return frame.V1
		}(),
		OutComponentID:     n.conf.OutComponentID,
		OutSignatureLinkID: linkID,
		OutKey:             n.conf.OutKey,
	})
	if err != nil {
		return nil, err
	}

	writeBuffer, err := ringbuffer.New(writeBufferSize)
	if err != nil {
		return nil, err
	}

	return &Channel{
		e:           e,
		label:       label,
		rwc:         rwc,
		n:           n,
		frw:         frw,
		writeBuffer: writeBuffer,
		terminate:   make(chan struct{}),
	}, nil
}

func (ch *Channel) close() {
	if ch.running {
		close(ch.terminate)
	} else {
		ch.rwc.Close()
	}
}

func (ch *Channel) start() {
	ch.running = true
	ch.n.channelsWg.Add(1)
	go ch.run()
}

func (ch *Channel) run() {
	defer ch.n.channelsWg.Done()

	readerDone := make(chan struct{})
	go ch.runReader(readerDone)

	writerDone := make(chan struct{})
	go ch.runWriter(writerDone)

	select {
	case <-readerDone:
		ch.rwc.Close()

		ch.writeBuffer.Close()
		<-writerDone

	case <-ch.terminate:
		ch.writeBuffer.Close()
		<-writerDone

		ch.rwc.Close()
		<-readerDone
	}

	ch.n.pushEvent(&EventChannelClose{ch})
	ch.n.closeChannel(ch)
}

func (ch *Channel) runReader(readerDone chan struct{}) {
	defer close(readerDone)

	// wait client here, in order to allow the writer goroutine to start
	// and allow clients to write messages before starting listening to events
	ch.n.pushEvent(&EventChannelOpen{ch})

	for {
		fr, err := ch.frw.Read()
		if err != nil {
			if _, ok := err.(*frame.ReadError); ok {
				ch.n.pushEvent(&EventParseError{err, ch})
				continue
			}
			return
		}

		evt := &EventFrame{fr, ch}

		if ch.n.nodeStreamRequest != nil {
			ch.n.nodeStreamRequest.onEventFrame(evt)
		}

		ch.n.pushEvent(evt)
	}
}

func (ch *Channel) runWriter(writerDone chan struct{}) {
	defer close(writerDone)

	for {
		what, ok := ch.writeBuffer.Pull()
		if !ok {
			return
		}

		switch wh := what.(type) {
		case message.Message:
			ch.frw.WriteMessage(wh) //nolint:errcheck

		case frame.Frame:
			ch.frw.WriteFrame(wh) //nolint:errcheck
		}
	}
}

// String implements fmt.Stringer.
func (ch *Channel) String() string {
	return ch.label
}

// Endpoint returns the channel Endpoint.
func (ch *Channel) Endpoint() Endpoint {
	return ch.e
}

func (ch *Channel) write(what interface{}) {
	ch.writeBuffer.Push(what)
}
