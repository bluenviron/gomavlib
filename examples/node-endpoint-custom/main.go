// Package main contains an example.
package main

import (
	"fmt"
	"log"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a custom endpoint from a io.ReadWriteCloser.
// 2) create a node which communicates with the custom endpoint.
// 3) print incoming messages.

// this is an example struct that implements io.ReadWriteCloser.
// it does not read anything and prints what it receives.
// the only requirement is that Close() must release Read().
type customEndpoint struct {
	readChan chan []byte
}

func newCustomEndpoint() *customEndpoint {
	return &customEndpoint{
		readChan: make(chan []byte),
	}
}

func (c *customEndpoint) Close() error {
	close(c.readChan)
	return nil
}

func (c *customEndpoint) Read(buf []byte) (int, error) {
	read, ok := <-c.readChan
	if !ok {
		return 0, fmt.Errorf("all right")
	}

	n := copy(buf, read)
	return n, nil
}

func (c *customEndpoint) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func main() {
	// allocate the custom endpoint
	endpoint := newCustomEndpoint()

	// create a node which communicates with the custom endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustom{ReadWriteCloser: endpoint},
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

	// queue a dummy message
	endpoint.readChan <- []byte{
		0xfd, 0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x02,
		0x03, 0x05, 0x03, 0xd9, 0xd1, 0x01, 0x02, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0e, 0x47, 0x04, 0x0c,
		0xef, 0x9b,
	}

	// print incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
