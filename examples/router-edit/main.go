package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v2"
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint
// 2) print incoming frames
// 3) edit messages of a specific kind
// 4) recompute the frame checksum and signature
// 5) route frame to every other channel

func main() {
	// create a node which communicates with a serial endpoint
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

	// print incoming frames
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())

			// if incoming message is a heartbeat
			if msg, ok := frm.Message().(*ardupilotmega.MessageHeartbeat); ok {
				// edit a field
				msg.Type = ardupilotmega.MAV_TYPE_SUBMARINE

				// since we changed the frame content, recompute checksum and signature
				err := node.FixFrame(frm.Frame)
				if err != nil {
					log.Printf("ERR: %v", err)
					continue
				}
			}

			// route frame to every other channel
			node.WriteFrameExcept(frm.Channel, frm.Frame)
		}
	}
}
