package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v2"
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint
// 2) print selected incoming messages

func main() {
	// create a node which communicates with a serial endpoint
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

	// print selected incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			switch msg := frm.Message().(type) {
			// if frm.Message() is a *ardupilotmega.MessageHeartbeat, access its fields
			case *ardupilotmega.MessageHeartbeat:
				log.Printf("received heartbeat (type %d)\n", msg.Type)

			// if frm.Message() is a *ardupilotmega.MessageServoOutputRaw, access its fields
			case *ardupilotmega.MessageServoOutputRaw:
				log.Printf("received servo output with values: %d %d %d %d %d %d %d %d\n",
					msg.Servo1Raw, msg.Servo2Raw, msg.Servo3Raw, msg.Servo4Raw,
					msg.Servo5Raw, msg.Servo6Raw, msg.Servo7Raw, msg.Servo8Raw)
			}
		}
	}
}
