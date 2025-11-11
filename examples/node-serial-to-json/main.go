// Package main contains an example.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
// 2) encode incoming messages into JSON.
// 3) print messages in the console.

func main() {
	// create a node which communicates with a serial endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			// encode incoming messages
			var enc []byte
			enc, err = json.Marshal(struct {
				Type    string
				Content any
			}{
				Type:    fmt.Sprintf("%T", frm.Message()),
				Content: filterFloats(frm.Message()),
			})
			if err != nil {
				panic(err)
			}

			// print messages in the console
			log.Printf("%s\n", enc)
		}
	}
}
