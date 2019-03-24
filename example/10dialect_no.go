// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
)

func main() {
	// create a node which
	// - communicates through a serial port
	// - does not use dialects
	// - writes messages with given system id and component id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:  nil,
		SystemId: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.NodeEventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())
		}
	}
}
