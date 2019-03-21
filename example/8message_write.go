// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	// create a node which
	// - communicates through a serial port
	// - understands ardupilotmega dialect
	// - writes messages with given system id and component id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.NodeEventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())

			// if message is a parameter read request addressed to this node
			if msg, ok := frm.Message().(*ardupilotmega.MessageParamRequestRead); ok &&
				msg.TargetSystem == 10 &&
				msg.TargetComponent == 1 &&
				msg.ParamId == "test_parameter" {

				// reply to sender (and no one else) by providing requested parameter
				node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageParamValue{
					ParamId:    "test_parameter",
					ParamValue: 123456,
					ParamType:  uint8(ardupilotmega.MAV_PARAM_TYPE_UINT32),
				})
			}
		}
	}
}
