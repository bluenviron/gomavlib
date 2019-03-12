// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
		Transports: []gomavlib.TransportConf{
			gomavlib.TransportUdpClient{"todo"},
		},
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for {
		frame, _, ok := node.ReadFrame()
		if ok == false {
			break
		}

		fmt.Printf("received: id=%d, %+v\n", frame.GetMessage().GetId(), frame.GetMessage())
	}
}
