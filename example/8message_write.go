// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	// create a node which understands given dialect, writes messages with given
	// system id and component id, and reads/writes through a serial port.
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
		Transports: []gomavlib.TransportConf{
			gomavlib.TransportSerial{"/dev/ttyAMA0", 57600},
		},
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// work in a loop
	for {
		// wait until a message is received.
		res, ok := node.Read()
		if ok == false {
			break
		}

		// print message details
		fmt.Printf("received: id=%d, %+v\n", res.Message().GetId(), res.Message())

		// if message is a parameter read request addressed to this node
		if msg, ok := res.Message().(*ardupilotmega.MessageParamRequestRead); ok &&
			msg.TargetSystem == 10 &&
			msg.TargetComponent == 1 &&
			msg.ParamId == "test_parameter" {

			// reply to sender (and no one else) by providing requested parameter
			node.WriteMessageTo(res.Channel(), &ardupilotmega.MessageParamValue{
				ParamId:    "test_parameter",
				ParamValue: 123456,
				ParamType:  uint8(ardupilotmega.MAV_PARAM_TYPE_UINT32),
			})
		}
	}
}
