package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
//    The node is configured to automatically request streams
//    to ardupilot devices, that require an explicit message.
// 2) print incoming frames.

func main() {
	// create a node
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:             ardupilotmega.Dialect,
		OutVersion:          gomavlib.V1, // Ardupilot uses V1
		OutSystemID:         10,
		StreamRequestEnable: true,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
