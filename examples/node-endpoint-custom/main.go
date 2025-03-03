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
	endpoint.readChan <- []byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b")

	// print incoming messages
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
