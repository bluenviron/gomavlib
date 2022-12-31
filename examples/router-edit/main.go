package main

import (
	"fmt"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
)

func main() {
	// create a node which
	// - communicates with multiple endpoints
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
			gomavlib.EndpointUDPClient{"1.2.3.4:5900"},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())

			// if incoming message is a heartbeat
			if msg, ok := frm.Message().(*ardupilotmega.MessageHeartbeat); ok {
				// edit a field of all incoming heartbeat messages
				msg.Type = ardupilotmega.MAV_TYPE_SUBMARINE

				// since we changed the frame content, recompute the frame checksum and signature
				err := node.FixFrame(frm.Frame)
				if err != nil {
					fmt.Printf("ERR: %v", err)
					continue
				}
			}
			// route frame to every other channel
			node.WriteFrameExcept(frm.Channel, frm.Frame)
		}
	}
}
