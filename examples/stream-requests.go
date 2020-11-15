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
	// - automatically requests streams to ardupilot devices
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:             ardupilotmega.Dialect,
		OutVersion:          gomavlib.V1, // Ardupilot uses V1
		OutSystemId:         10,
		StreamRequestEnable: true,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())
		}
	}
}
