package main

import (
	"fmt"
	"log"

	"github.com/aler9/gomavlib"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
)

// this is an example struct that implements io.ReadWriteCloser.
// it does not read anything and prints what it receives.
// the only requirement is that Close() must release Read().
type CustomEndpoint struct {
	readChan chan []byte
}

func NewCustomEndpoint() *CustomEndpoint {
	return &CustomEndpoint{
		readChan: make(chan []byte),
	}
}

func (c *CustomEndpoint) Close() error {
	close(c.readChan)
	return nil
}

func (c *CustomEndpoint) Read(buf []byte) (int, error) {
	read, ok := <-c.readChan
	if !ok {
		return 0, fmt.Errorf("all right")
	}

	n := copy(buf, read)
	return n, nil
}

func (c *CustomEndpoint) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func main() {
	// allocate the custom endpoint
	endpoint := NewCustomEndpoint()

	// create a node which
	// - communicates with a custom endpoint
	// - understands ardupilotmega dialect
	// - writes messages with given system id
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustom{endpoint},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// queue a dummy message
	endpoint.readChan <- []byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b")

	// print every message we receive
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
		}
	}
}
