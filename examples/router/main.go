package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v2"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint
// 2) print incoming frames
// 3) route received frames to every other channel

func main() {
	// create a node which communicates with a serial endpoint
	// - communicates with multiple endpoints
	// -
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
			gomavlib.EndpointUDPClient{"1.2.3.4:5900"},
		},
		Dialect:     nil,         // do not use a dialect and do not attempt to decode messages (in a router it is preferable)
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print incoming frames
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())

			// route frame to every other channel
			node.WriteFrameExcept(frm.Channel, frm.Frame)
		}
	}
}
