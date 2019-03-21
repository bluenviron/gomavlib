package gomavlib

// NodeEvent is the interface implemented by all events received through node.Events()
type NodeEvent interface {
	isEvent()
}

// NodeEventFrame is the event fired when a frame is received
type NodeEventFrame struct {
	// the frame
	Frame Frame

	// a parse error returned instead of Frame
	// This is used only when ReturnParseErrors is true
	Error error

	// the channel used to send the frame
	Channel *EndpointChannel
}

func (*NodeEventFrame) isEvent() {}
