// Package main contains an example.
package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a custom TCP/TLS endpoint in client mode.
// 2) print incoming messages.

func main() {
	// create a node which communicates with a TCP endpoint in client mode
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustomClient{
				Address: "127.0.0.1:5600",
				Connect: func(address string) (net.Conn, error) {
					tlsConfig := &tls.Config{
						// skip checking the certificate against a CA (just set to true for simplicity of this example)
						InsecureSkipVerify: true,
					}

					return tls.Dial("tcp", address, tlsConfig)
				},
				Label: "TCP/TLS",
			},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
