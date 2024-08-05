package gomavlib

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

const (
	writeBufferSize = 64
)

type broadcastMessage struct {
	message.Message
}

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
	n     *Node
	e     Endpoint
	label string
	rwc   io.Closer

	ctx       context.Context
	ctxCancel func()
	frw       *frame.ReadWriter
	running   bool

	// in
	chWrite chan interface{}

	// out
	done chan struct{}
}

func newChannel(
	n *Node,
	e Endpoint,
	label string,
	rwc io.ReadWriteCloser,
) (*Channel, error) {
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

	ctx, ctxCancel := context.WithCancel(context.Background())

	return &Channel{
		n:         n,
		e:         e,
		label:     label,
		rwc:       rwc,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		frw:       frw,
		chWrite:   make(chan interface{}, writeBufferSize),
		done:      make(chan struct{}),
	}, nil
}

func (ch *Channel) close() {
	ch.ctxCancel()
	if !ch.running {
		ch.rwc.Close()
	}
}

func (ch *Channel) start() {
	ch.running = true
	ch.n.wg.Add(1)
	go ch.run()
}

func (ch *Channel) run() {
	defer close(ch.done)
	defer ch.n.wg.Done()

	readerDone := make(chan struct{})
	go ch.runReader(readerDone)

	writerTerminate := make(chan struct{})
	writerDone := make(chan struct{})
	go ch.runWriter(writerTerminate, writerDone)

	select {
	case <-readerDone:
		ch.rwc.Close()

		close(writerTerminate)
		<-writerDone

	case <-ch.ctx.Done():
		close(writerTerminate)
		<-writerDone

		ch.rwc.Close()
		<-readerDone
	}

	ch.ctxCancel()

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
			var eerr frame.ReadError
			if errors.As(err, &eerr) {
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

func (ch *Channel) runWriter(writerTerminate chan struct{}, writerDone chan struct{}) {
	defer close(writerDone)

	for {
		select {
		case what := <-ch.chWrite:
			switch what := what.(type) {
			case broadcastMessage:
				ch.frw.WriteBroadcastMessage(what.Message) //nolint:errcheck

			case message.Message:
				ch.frw.WriteMessage(what) //nolint:errcheck

			case frame.Frame:
				ch.frw.WriteFrame(what) //nolint:errcheck
			}

		case <-writerTerminate:
			return
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
	select {
	case ch.chWrite <- what:
	case <-ch.ctx.Done():
	default: // buffer is full
	}
}
