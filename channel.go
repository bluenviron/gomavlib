package gomavlib

import (
	"io"
)

// Channel is a communication channel created by an endpoint. For instance, a
// TCP client endpoint creates a single channel, while a TCP server endpoint
// creates a channel for each incoming connection.
type Channel struct {
	// the endpoint which the channel belongs to
	Endpoint Endpoint

	label  string
	rwc    io.ReadWriteCloser
	n      *Node
	parser *Parser

	writeChan  chan interface{}
	allWritten chan struct{}
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) (*Channel, error) {
	parser, err := NewParser(ParserConf{
		Reader:             rwc,
		Writer:             rwc,
		Dialect:            n.conf.Dialect,
		InKey:              n.conf.InKey,
		OutSystemId:        n.conf.OutSystemId,
		OutVersion:         n.conf.OutVersion,
		OutComponentId:     n.conf.OutComponentId,
		OutSignatureLinkId: randomByte(),
		OutKey:             n.conf.OutKey,
	})
	if err != nil {
		return nil, err
	}

	return &Channel{
		Endpoint:   e,
		label:      label,
		rwc:        rwc,
		n:          n,
		parser:     parser,
		writeChan:  make(chan interface{}),
		allWritten: make(chan struct{}),
	}, nil
}

// String implements fmt.Stringer and returns the channel label.
func (ch *Channel) String() string {
	return ch.label
}

func (ch *Channel) close() {
	// wait until all frame have been written
	close(ch.writeChan)
	<-ch.allWritten

	// close reader/writer after ensuring all frames have been written
	ch.rwc.Close()
}

func (ch *Channel) run() {
	// reader
	readerDone := make(chan struct{})
	go func() {
		defer close(readerDone)
		defer func() { ch.n.eventsOut <- &EventChannelClose{ch} }()
		defer func() { ch.n.eventsIn <- &eventInChannelClosed{ch} }()

		ch.n.eventsOut <- &EventChannelOpen{ch}

		for {
			frame, err := ch.parser.Read()
			if err != nil {
				// continue in case of parse errors
				if _, ok := err.(*ParserError); ok {
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

	// writer
	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		defer func() { ch.allWritten <- struct{}{} }()

		for what := range ch.writeChan {
			switch wh := what.(type) {
			case Message:
				ch.parser.WriteMessage(wh)

			case Frame:
				ch.parser.WriteFrame(wh)
			}
		}
	}()

	<-readerDone
	<-writerDone
}
