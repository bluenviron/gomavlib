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
	// - communicates with a serial port.
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	// - validates incoming messages via InSignatureKey
	// - sign outgoing messages via OutSignatureKey
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{"/dev/ttyUSB0:57600"},
		},
		Dialect:         ardupilotmega.Dialect,
		OutSystemId:     10,
		InSignatureKey:  key,
		OutSignatureKey: key,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetId(), frm.Message())
		}
	}
}
