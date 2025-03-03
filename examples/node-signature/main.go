package main

import (
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
//    The node is configured to validate incoming messages through InKey
//    and sign outgoing messages with OutKey.
// 2) print incoming frames.

func main() {
	// initialize a 6-bytes key. A key can have up to 32 bytes.
	key := frame.NewV2Key([]byte("abcdef"))

	// create a node.
	node := &gomavlib.Node{
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
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
