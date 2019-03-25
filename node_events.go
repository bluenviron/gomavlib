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

// Message returns the message inside the frame.
func (res *NodeEventFrame) Message() Message {
	return res.Frame.GetMessage()
}

// SystemId returns the frame system id.
func (res *NodeEventFrame) SystemId() byte {
	return res.Frame.GetSystemId()
}

// ComponentId returns the frame component id.
func (res *NodeEventFrame) ComponentId() byte {
	return res.Frame.GetComponentId()
}

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
	Channel *EndpointChannel
}

func (*NodeEventChannelOpen) isEvent() {}

// NodeEventChannelClose is the event fired when a channel is closed
type NodeEventChannelClose struct {
	Channel *EndpointChannel
}

func (*NodeEventChannelClose) isEvent() {}

// NodeEventNodeAppear is the event fired when a new node is detected
type NodeEventNodeAppear struct {
	Node NodeIndex
}

func (*NodeEventNodeAppear) isEvent() {}

// NodeEventNodeDisappear is the event fired when a node disappears (i.e. times out)
type NodeEventNodeDisappear struct {
	Node NodeIndex
}

func (*NodeEventNodeDisappear) isEvent() {}
