package gomavlib

type eventIn interface {
	isEventIn()
}

type eventInChannelNew struct {
	ch *Channel
}

func (*eventInChannelNew) isEventIn() {}

type eventInChannelClosed struct {
	ch *Channel
}

func (*eventInChannelClosed) isEventIn() {}

type eventInWriteTo struct {
	ch   *Channel
	what interface{}
}

func (*eventInWriteTo) isEventIn() {}

type eventInWriteAll struct {
	what interface{}
}

func (*eventInWriteAll) isEventIn() {}

type eventInWriteExcept struct {
	except *Channel
	what   interface{}
}

func (*eventInWriteExcept) isEventIn() {}

type eventInClose struct {
}

func (*eventInClose) isEventIn() {}

// Event is the interface implemented by all events received through node.Events().
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
	Frame Frame

	// the channel from which the frame was received
	Channel *Channel
}

func (*EventFrame) isEventOut() {}

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

func (*EventParseError) isEventOut() {}

// EventStreamRequested is the event fired when an automatic stream request is sent.
type EventStreamRequested struct {
	// the channel to which the stream request is addressed
	Channel *Channel
	// the system id to which the stream requests is addressed
	SystemId byte
	// the component id to which the stream requests is addressed
	ComponentId byte
}

func (*EventStreamRequested) isEventOut() {}
