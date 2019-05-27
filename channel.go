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

	label      string
	rwc        io.ReadWriteCloser
	n          *Node
	parser     *Parser
	writeChan  chan interface{}
	allWritten chan struct{}
}

// String implements fmt.Stringer and returns the channel label.
func (e *Channel) String() string {
	return e.label
}

func newChannel(n *Node, e Endpoint, label string, rwc io.ReadWriteCloser) *Channel {
	parser, _ := NewParser(ParserConf{
		Reader:             rwc,
		Writer:             rwc,
		Dialect:            n.conf.Dialect,
		InSignatureKey:     n.conf.InSignatureKey,
		OutSystemId:        n.conf.OutSystemId,
		OutComponentId:     n.conf.OutComponentId,
		OutSignatureLinkId: randomByte(),
		OutSignatureKey:    n.conf.OutSignatureKey,
	})

	ch := &Channel{
		Endpoint:   e,
		label:      label,
		rwc:        rwc,
		n:          n,
		parser:     parser,
		writeChan:  make(chan interface{}),
		allWritten: make(chan struct{}),
	}

	return ch
}

func (ch *Channel) close() {
	// wait until all packets have been written
	close(ch.writeChan)
	<-ch.allWritten
}

func (ch *Channel) run() {
	// reader
	readerDone := make(chan struct{})
	go func() {
		defer func() { readerDone <- struct{}{} }()
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
		defer func() { writerDone <- struct{}{} }()
		defer func() { ch.allWritten <- struct{}{} }()
		defer ch.rwc.Close()

		for what := range ch.writeChan {
			switch wh := what.(type) {
			case Message:
				if ch.n.conf.OutVersion == V1 {
					ch.parser.Write(&FrameV1{Message: wh}, false)
				} else {
					ch.parser.Write(&FrameV2{Message: wh}, false)
				}

			case Frame:
				ch.parser.Write(wh, true)
			}
		}
	}()

	<-readerDone
	<-writerDone
}
