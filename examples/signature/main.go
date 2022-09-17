package main

import (
	"fmt"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
	"github.com/aler9/gomavlib/pkg/frame"
)

func main() {
	// initialize a 6-bytes key. A key can have up to 32 bytes.
	key := frame.NewV2Key([]byte("abcdef"))

	// create a node which
	// - communicates with a serial port.
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	// - validates incoming messages via InKey
	// - sign outgoing messages via OutKey
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // V2 is mandatory for signatures
		OutSystemID: 10,
		InKey:       key,
		OutKey:      key,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
