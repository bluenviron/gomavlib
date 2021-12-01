package gomavlib

import (
	"io"
	"math/rand"

	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/aler9/gomavlib/pkg/parser"
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
	reader  *parser.Reader
	writer  *parser.Writer
	running bool

	// in
	write     chan interface{}
	terminate chan struct{}
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) (*Channel, error) {
	reader, err := parser.NewReader(parser.ReaderConf{
		Reader:    rwc,
		DialectDE: n.dialectDE,
		InKey:     n.conf.InKey,
	})
	if err != nil {
		return nil, err
	}

	writer, err := parser.NewWriter(parser.WriterConf{
		Writer:      rwc,
		DialectDE:   n.dialectDE,
		OutSystemID: n.conf.OutSystemID,
		OutVersion: func() parser.WriterOutVersion {
			if n.conf.OutVersion == V2 {
				return parser.V2
			}
			return parser.V1
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
		reader:    reader,
		writer:    writer,
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
		ch.n.events <- &EventChannelOpen{ch}

		for {
			frame, err := ch.reader.Read()
			if err != nil {
				// continue in case of parse errors
				if _, ok := err.(*parser.ReadError); ok {
					ch.n.events <- &EventParseError{err, ch}
					continue
				}
				return
			}

			evt := &EventFrame{frame, ch}

			if ch.n.nodeStreamRequest != nil {
				ch.n.nodeStreamRequest.onEventFrame(evt)
			}

			ch.n.events <- evt
		}
	}()

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)

		for what := range ch.write {
			switch wh := what.(type) {
			case msg.Message:
				ch.writer.WriteMessage(wh)

			case frame.Frame:
				ch.writer.WriteFrame(wh)
			}
		}
	}()

	select {
	case <-readerDone:
		ch.n.events <- &EventChannelClose{ch}

		ch.n.channelClose <- ch
		<-ch.terminate

		close(ch.write)
		<-writerDone

		ch.rwc.Close()

	case <-ch.terminate:
		ch.n.events <- &EventChannelClose{ch}

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
