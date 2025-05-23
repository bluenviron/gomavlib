// Package main contains an example.
package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// this example shows how to:
// 1) create a custom dialect from a list of messages.
// 2) create a node which understands the custom dialect.
// 3) print incoming messages.

// MessageCustom is a custom message.
// It must be prefixed with "Message" and implement the message.Message interface.
type MessageCustom struct {
	Param1 uint8
	Param2 uint8
	Param3 uint32
}

// GetID implements the message.Message interface.
func (*MessageCustom) GetID() uint32 {
	return 304
}

func main() {
	// create a custom dialect from a list of messages
	dialect := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageCustom{},
		},
	}

	// create a node which understands the custom dialect
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
