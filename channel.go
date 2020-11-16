package gomavlib

import (
	"io"

	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/aler9/gomavlib/pkg/transceiver"
)

// Channel is a communication channel created by an Endpoint.
// An Endpoint can create channels.
// For instance, a TCP client endpoint creates a single channel, while a TCP
// server endpoint creates a channel for each incoming connection.
type Channel struct {
	// the endpoint which the channel belongs to
	Endpoint Endpoint

	label       string
	rwc         io.ReadWriteCloser
	n           *Node
	transceiver *transceiver.Transceiver

	write     chan interface{}
	terminate chan struct{}
	done      chan struct{}
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) (*Channel, error) {
	transceiver, err := transceiver.New(transceiver.TransceiverConf{
		Reader:      rwc,
		Writer:      rwc,
		DialectDE:   n.dialectDE,
		InKey:       n.conf.InKey,
		OutSystemId: n.conf.OutSystemId,
		OutVersion: func() transceiver.Version {
			if n.conf.OutVersion == V2 {
				return transceiver.V2
			}
			return transceiver.V1
		}(),
		OutComponentId:     n.conf.OutComponentId,
		OutSignatureLinkId: randomByte(),
		OutKey:             n.conf.OutKey,
	})
	if err != nil {
		return nil, err
	}

	return &Channel{
		Endpoint:    e,
		label:       label,
		rwc:         rwc,
		n:           n,
		transceiver: transceiver,
		write:       make(chan interface{}),
		terminate:   make(chan struct{}),
		done:        make(chan struct{}),
	}, nil
}

// String implements fmt.Stringer and returns the channel label.
func (ch *Channel) String() string {
	return ch.label
}

func (ch *Channel) run() {
	defer close(ch.done)

	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)

		// wait client here, in order to allow the writer goroutine to start
		// and allow clients to write messages before starting listening to events
		ch.n.eventsOut <- &EventChannelOpen{ch}

		for {
			frame, err := ch.transceiver.Read()
			if err != nil {
				// continue in case of parse errors
				if _, ok := err.(*transceiver.TransceiverError); ok {
					ch.n.eventsOut <- &EventParseError{err, ch}
					continue
				}
				return
			}

			evt := &EventFrame{frame, ch}

			if ch.n.nodeStreamRequest != nil {
				ch.n.nodeStreamRequest.onEventFrame(evt)
			}

			ch.n.eventsOut <- evt
		}
	}()

	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)

		for what := range ch.write {
			switch wh := what.(type) {
			case msg.Message:
				ch.transceiver.WriteMessage(wh)

			case frame.Frame:
				ch.transceiver.WriteFrame(wh)
			}
		}
	}()

	select {
	case <-readerDone:
		ch.n.eventsOut <- &EventChannelClose{ch}

		ch.n.channelClose <- ch
		<-ch.terminate

		close(ch.write)
		<-writerDone

		ch.rwc.Close()

	case <-ch.terminate:
		ch.n.eventsOut <- &EventChannelClose{ch}

		close(ch.write)
		<-writerDone

		ch.rwc.Close()
		<-readerDone
	}
}
