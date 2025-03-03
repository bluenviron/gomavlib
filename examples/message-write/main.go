package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
// 2) wait for a specific incoming message.
// 3) write a reply message.

func main() {
	// create a node which communicates with a serial endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())

			// if message is a parameter read request addressed to this node
			if msg, ok := frm.Message().(*ardupilotmega.MessageParamRequestRead); ok &&
				msg.TargetSystem == 10 &&
				msg.TargetComponent == 1 &&
				msg.ParamId == "test_parameter" {

				// reply to sender (and no one else) and provide the requested parameter
				node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageParamValue{
					ParamId:    "test_parameter",
					ParamValue: 123456,
					ParamType:  ardupilotmega.MAV_PARAM_TYPE_UINT32,
				})
			}
		}
	}
}
