package gomavlib

// NodeEvent is the interface implemented by all events received through node.Events()
type NodeEvent interface {
	isEvent()
}

// NodeEventFrame is the event fired when a frame is received
type NodeEventFrame struct {
	// the frame
	Frame Frame

	// the channel used to send the frame
	Channel *EndpointChannel
}

func (*NodeEventFrame) isEvent() {}

// NodeEventParseError is the event fired when a parse error occurs
type NodeEventParseError struct {
	// the error
	Error error

	// the channel used to send the frame
	Channel *EndpointChannel
}

func (*NodeEventParseError) isEvent() {}

// NodeEventChannelOpen is the event fired when a channel is opened
type NodeEventChannelOpen struct {
	// the channel
	Channel *EndpointChannel
}

func (*NodeEventChannelOpen) isEvent() {}

// NodeEventChannelClose is the event fired when a channel is closed
type NodeEventChannelClose struct {
	// the channel
	Channel *EndpointChannel
}

func (*NodeEventChannelClose) isEvent() {}
