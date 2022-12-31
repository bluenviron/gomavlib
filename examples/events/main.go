package main

import (
	"log"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
)

func main() {
	// create a node which
	// - communicates with a serial port
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// gomavlib provides different kinds of event
	for evt := range node.Events() {
		switch ee := evt.(type) {
		case *gomavlib.EventFrame:
			log.Printf("frame received: %v\n", ee)

		case *gomavlib.EventParseError:
			log.Printf("parse error: %v\n", ee)

		case *gomavlib.EventChannelOpen:
			log.Printf("channel opened: %v\n", ee)

		case *gomavlib.EventChannelClose:
			log.Printf("channel closed: %v\n", ee)
		}
	}
}
