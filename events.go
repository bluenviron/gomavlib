package gomavlib

import (
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
)

// Event is the interface implemented by all events received with node.Events().
type Event interface {
	isEventOut()
}

// EventChannelOpen is the event fired when a channel gets opened.
type EventChannelOpen struct {
	Channel *Channel
}

func (*EventChannelOpen) isEventOut() {}

// EventChannelClose is the event fired when a channel gets closed.
type EventChannelClose struct {
	Channel *Channel
}

func (*EventChannelClose) isEventOut() {}

// EventFrame is the event fired when a frame is received.
type EventFrame struct {
	// the frame
	Frame frame.Frame

	// the channel from which the frame was received
	Channel *Channel
}

func (*EventFrame) isEventOut() {}

// SystemID returns the frame system id.
func (res *EventFrame) SystemID() byte {
	return res.Frame.GetSystemID()
}

// ComponentID returns the frame component id.
func (res *EventFrame) ComponentID() byte {
	return res.Frame.GetComponentID()
}

// Message returns the message inside the frame.
func (res *EventFrame) Message() msg.Message {
	return res.Frame.GetMessage()
}

// EventParseError is the event fired when a parse error occurs.
type EventParseError struct {
	// the error
	Error error

	// the channel used to send the frame
	Channel *Channel
}

func (*EventParseError) isEventOut() {}

// EventStreamRequested is the event fired when an automatic stream request is sent.
type EventStreamRequested struct {
	// the channel to which the stream request is addressed
	Channel *Channel
	// the system id to which the stream requests is addressed
	SystemID byte
	// the component id to which the stream requests is addressed
	ComponentID byte
}

func (*EventStreamRequested) isEventOut() {}
