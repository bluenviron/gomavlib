package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint
// 2) convert any message to a generic JSON object
// the output of this example can be combined with a JSON parser like jq to display individual messages, for example:
//     go run ./examples/serial-to-json | jq 'select(.MessageType == "*minimal.MessageHeartbeat")'

func main() {
	baudRate := flag.Int("b", 57600, "baud rate")
	port := flag.String("p", "/dev/ttyUSB0", "port")
	flag.Parse()

	// create a node which communicates with a serial endpoint
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: *port,
				Baud:   *baudRate,
			},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	jsonEncoder := json.NewEncoder(os.Stdout)
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			// add the message type to the JSON object
			if err := jsonEncoder.Encode(struct {
				MessageType string
				Frame       *frame.Frame
			}{
				MessageType: fmt.Sprintf("%T", frm.Message()),
				Frame:       &frm.Frame,
			}); err != nil {
				// some messages contain floating point NaN values which cannot be encoded in JSON
				// silently ignore these errors
				if err.Error() != "json: unsupported value: NaN" {
					panic(err)
				}
			}
		}
	}
}
