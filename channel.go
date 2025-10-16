package gomavlib

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

const (
	writeBufferSize = 64
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
	node     *Node
	endpoint Endpoint
	label    string
	rwc      io.ReadWriteCloser

	ctx          context.Context
	ctxCancel    func()
	frameWriter  *frame.ReadWriter
	streamWriter *streamwriter.Writer
	running      bool

	// in
	chWrite chan interface{}

	// out
	done chan struct{}
}

func (ch *Channel) initialize() error {
	linkID, err := randomByte()
	if err != nil {
		return err
	}

	ch.frameWriter = &frame.ReadWriter{
		ByteReadWriter: ch.rwc,
		DialectRW:      ch.node.dialectRW,
		InKey:          ch.node.InKey,
	}
	err = ch.frameWriter.Initialize()
	if err != nil {
		return err
	}

	ch.streamWriter = &streamwriter.Writer{
		FrameWriter: ch.frameWriter.Writer,
		SystemID:    ch.node.OutSystemID,
		Version: func() streamwriter.Version {
			if ch.node.OutVersion == V2 {
				return streamwriter.V2
			}
			return streamwriter.V1
		}(),
		ComponentID:     ch.node.OutComponentID,
		SignatureLinkID: linkID,
		Key:             ch.node.OutKey,
	}
	err = ch.streamWriter.Initialize()
	if err != nil {
		return err
	}

	ch.ctx, ch.ctxCancel = context.WithCancel(context.Background())
	ch.chWrite = make(chan interface{}, writeBufferSize)
	ch.done = make(chan struct{})

	return nil
}

func (ch *Channel) close() {
	ch.ctxCancel()
	if !ch.running {
		ch.rwc.Close()
	}
}

func (ch *Channel) start() {
	ch.running = true
	ch.node.wg.Add(1)
	go ch.run()
}

func (ch *Channel) run() {
	defer close(ch.done)
	defer ch.node.wg.Done()

	readerDone := make(chan error)
	go func() {
		readerDone <- ch.runReader()
	}()

	writerTerminate := make(chan struct{})
	writerDone := make(chan error)
	go func() {
		writerDone <- ch.runWriter(writerTerminate)
	}()

	var err error

	select {
	case err = <-readerDone:
		ch.rwc.Close()

		close(writerTerminate)
		<-writerDone

	case err = <-writerDone:
		ch.rwc.Close()
		<-readerDone

	case <-ch.ctx.Done():
		close(writerTerminate)
		<-writerDone

		ch.rwc.Close()
		<-readerDone
	}

	ch.ctxCancel()

	ch.node.pushEvent(&EventChannelClose{
		Channel: ch,
		Error:   err,
	})
	ch.node.closeChannel(ch)
}

func (ch *Channel) runReader() error {
	// wait client here, in order to allow the writer goroutine to start
	// and allow clients to write messages before starting listening to events
	ch.node.pushEvent(&EventChannelOpen{ch})

	for {
		fr, err := ch.frameWriter.Read()
		if err != nil {
			var eerr frame.ReadError
			if errors.As(err, &eerr) {
				ch.node.pushEvent(&EventParseError{err, ch})
				continue
			}
			return err
		}

		evt := &EventFrame{fr, ch}

		if ch.node.nodeStreamRequest != nil {
			ch.node.nodeStreamRequest.onEventFrame(evt)
		}

		ch.node.pushEvent(evt)
	}
}

func (ch *Channel) runWriter(writerTerminate chan struct{}) error {
	for {
		select {
		case what := <-ch.chWrite:
			switch wh := what.(type) {
			case message.Message:
				err := ch.streamWriter.Write(wh)
				if err != nil {
					return err
				}

			case frame.Frame:
				err := ch.frameWriter.Write(wh)
				if err != nil {
					return err
				}
			}

		case <-writerTerminate:
			return nil
		}
	}
}

// String implements fmt.Stringer.
func (ch *Channel) String() string {
	return ch.label
}

// Endpoint returns the channel Endpoint.
func (ch *Channel) Endpoint() Endpoint {
	return ch.endpoint
}

func (ch *Channel) write(what interface{}) {
	select {
	case ch.chWrite <- what:
	case <-ch.ctx.Done():
	default: // buffer is full
	}
}
