package gomavlib

// Event is the interface implemented by all events received through node.Events().
type Event interface {
	isEvent()
}

// EventChannelOpen is the event fired when a channel is opened.
type EventChannelOpen struct {
	Channel *Channel
}

func (*EventChannelOpen) isEvent() {}

// EventChannelClose is the event fired when a channel is closed.
type EventChannelClose struct {
	Channel *Channel
}

func (*EventChannelClose) isEvent() {}

// EventFrame is the event fired when a frame is received.
type EventFrame struct {
	// the frame
	Frame Frame

	// the channel from which the frame was received
	Channel *Channel
}

func (*EventFrame) isEvent() {}

// SystemId returns the frame system id.
func (res *EventFrame) SystemId() byte {
	return res.Frame.GetSystemId()
}

// ComponentId returns the frame component id.
func (res *EventFrame) ComponentId() byte {
	return res.Frame.GetComponentId()
}

// Message returns the message inside the frame.
func (res *EventFrame) Message() Message {
	return res.Frame.GetMessage()
}

// EventParseError is the event fired when a parse error occurs.
type EventParseError struct {
	// the error
	Error error

	// the channel used to send the frame
	Channel *Channel
}

func (*EventParseError) isEvent() {}
