// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	// create a node which understands given dialect, writes messages with given
	// system id and component id, and reads/writes to an UDP endpoint in broadcast mode.
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUdpBroadcast{BroadcastAddress: "192.168.7.255:5600"},
		},
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for {
		// wait until a message is received.
		res, ok := node.Read()
		if ok == false {
			break
		}

		// print message details
		fmt.Printf("received: id=%d, %+v\n", res.Message().GetId(), res.Message())
	}
}
