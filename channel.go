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

	label     string
	rwc       io.ReadWriteCloser
	writeChan chan interface{}
	n         *Node
	parser    *Parser
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
		Endpoint:  e,
		label:     label,
		rwc:       rwc,
		writeChan: make(chan interface{}),
		n:         n,
		parser:    parser,
	}

	return ch
}

func (ch *Channel) close() {
	ch.rwc.Close()
}

func (ch *Channel) start() {
	ch.n.wg.Add(2)

	// reader
	go func() {
		defer ch.n.wg.Done()

		defer func() {
			ch.n.channelsMutex.Lock()
			delete(ch.n.channels, ch)
			ch.n.channelsMutex.Unlock()
			close(ch.writeChan)
			ch.n.eventChan <- &EventChannelClose{ch}
		}()

		ch.n.eventChan <- &EventChannelOpen{ch}

		for {
			frame, err := ch.parser.Read()
			if err != nil {
				// continue in case of parse errors
				if _, ok := err.(*ParserError); ok {
					ch.n.eventChan <- &EventParseError{err, ch}
					continue
				}
				// avoid calling twice Close()
				if err != errorTerminated {
					ch.rwc.Close()
				}
				return
			}

			ch.n.eventChan <- &EventFrame{frame, ch}
		}
	}()

	// writer
	go func() {
		defer ch.n.wg.Done()

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

			ch.n.writeDone <- struct{}{}
		}
	}()
}
