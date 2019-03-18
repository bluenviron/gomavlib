// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	// initialize a 6-bytes key. A key can have up to 32 bytes.
	key := gomavlib.NewFrameSignatureKey([]byte("abcdef"))

	// create a node which
	// - communicates through a serial port.
	// - understands ardupilotmega dialect
	// - writes messages with given system id and component id
	// - validates incoming messages via SignatureInKey
	// - sign outgoing messages via SignatureOutKey
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:         ardupilotmega.Dialect,
		SystemId:        10,
		ComponentId:     1,
		SignatureInKey:  key,
		SignatureOutKey: key,
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
