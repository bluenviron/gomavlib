package gomavlib

import (
	"crypto/rand"
	"io"

	"github.com/bluenviron/gomavlib/v2/pkg/frame"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

func randomByte() byte {
	var buf [1]byte
	rand.Read(buf[:])
	return buf[0]
}

// Channel is a communication channel created by an Endpoint.
// An Endpoint can create channels.
// For instance, a TCP client endpoint creates a single channel, while a TCP
// server endpoint creates a channel for each incoming connection.
type Channel struct {
	e       Endpoint
	label   string
	rwc     io.ReadWriteCloser
	n       *Node
	frw     *frame.ReadWriter
	running bool

	// in
	write     chan interface{}
	terminate chan struct{}
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) (*Channel, error) {
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
		OutSignatureLinkID: randomByte(),
		OutKey:             n.conf.OutKey,
	})
	if err != nil {
		return nil, err
	}

	return &Channel{
		e:         e,
		label:     label,
		rwc:       rwc,
		n:         n,
		frw:       frw,
		write:     make(chan interface{}),
		terminate: make(chan struct{}),
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
	go func() {
		defer close(readerDone)

		// wait client here, in order to allow the writer goroutine to start
		// and allow clients to write messages before starting listening to events
		select {
		case ch.n.events <- &EventChannelOpen{ch}:
		case <-ch.n.terminate:
		}

		for {
			fr, err := ch.frw.Read()
			if err != nil {
				// ignore parse errors
				if _, ok := err.(*frame.ReadError); ok {
					select {
					case ch.n.events <- &EventParseError{err, ch}:
					case <-ch.n.terminate:
					}
					continue
				}
				return
			}

			evt := &EventFrame{fr, ch}

			if ch.n.nodeStreamRequest != nil {
				ch.n.nodeStreamRequest.onEventFrame(evt)
			}

			select {
			case ch.n.events <- evt:
			case <-ch.n.terminate:
			}
		}
	}()

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)

		for what := range ch.write {
			switch wh := what.(type) {
			case message.Message:
				ch.frw.WriteMessage(wh)

			case frame.Frame:
				ch.frw.WriteFrame(wh)
			}
		}
	}()

	select {
	case <-readerDone:
		select {
		case ch.n.events <- &EventChannelClose{ch}:
		case <-ch.n.terminate:
		}

		select {
		case ch.n.channelClose <- ch:
		case <-ch.n.terminate:
		}

		<-ch.terminate

		close(ch.write)
		<-writerDone

		ch.rwc.Close()

	case <-ch.terminate:
		select {
		case ch.n.events <- &EventChannelClose{ch}:
		case <-ch.n.terminate:
		}

		close(ch.write)
		<-writerDone

		ch.rwc.Close()
		<-readerDone
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
