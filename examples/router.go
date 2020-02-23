// +build ignore

package main

import (
	"fmt"

	"github.com/aler9/gomavlib"
)

func main() {
	// create a node which
	// - communicates with multiple endpoints
	// - is dialect agnostic, does not attempt to decode messages (in a router it is preferable)
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
			gomavlib.EndpointUdpClient{"1.2.3.4:5900"},
		},
		Dialect:     nil,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to write to the target
		OutSystemId: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())

			// route frame to every other channel
			node.WriteFrameExcept(frm.Channel, frm.Frame)
		}
	}
}
