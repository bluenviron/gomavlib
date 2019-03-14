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

	// create a node which understands given dialect, writes messages with given
	// system id and component id, and reads/writes through a serial port.
	// incoming messages are verified via SignatureInKey, and outgoing messages
	// are signed via SignatureOutKey.
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
		Transports: []gomavlib.TransportConf{
			gomavlib.TransportSerial{"/dev/ttyAMA0", 57600},
		},
		SignatureInKey:  key,
		SignatureOutKey: key,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// work in a loop
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
