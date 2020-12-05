// +build ignore

package main

import (
	"fmt"

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
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print selected messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {

			switch msg := frm.Message().(type) {
			// if frm.Message() is a *ardupilotmega.MessageHeartbeat, access its fields
			case *ardupilotmega.MessageHeartbeat:
				fmt.Printf("received heartbeat (type %d)\n", msg.Type)

			// if frm.Message() is a *ardupilotmega.MessageServoOutputRaw, access its fields
			case *ardupilotmega.MessageServoOutputRaw:
				fmt.Printf("received servo output with values: %d %d %d %d %d %d %d %d\n",
					msg.Servo1Raw, msg.Servo2Raw, msg.Servo3Raw, msg.Servo4Raw,
					msg.Servo5Raw, msg.Servo6Raw, msg.Servo7Raw, msg.Servo8Raw)
			}
		}
	}
}
