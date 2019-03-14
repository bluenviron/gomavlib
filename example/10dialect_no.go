// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
)

func main() {
	// create a node which does not use dialects, writes messages with given
	// system id and component id, and reads/writes through a serial port.
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Dialect:     nil,
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
