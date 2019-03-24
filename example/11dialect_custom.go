// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
)

// this is a custom message.
// It must be prefixed with Message and implement the gomavlib.Message interface.
type MessageMyCustomMessage struct {
	Param1 uint8
	Param2 uint8
	Param3 uint32
}

func (*MessageMyCustomMessage) GetId() uint32 {
	return 304
}

func main() {
	// create a custom dialect from messages
	dialect, err := gomavlib.NewDialect([]gomavlib.Message{
		&MessageMyCustomMessage{},
	})
	if err != nil {
		panic(err)
	}

	// create a node which
	// - communicates through a serial port
	// - understands our custom dialect
	// - writes messages with given system id and component id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:  dialect,
		SystemId: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.NodeEventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())
		}
	}
}
