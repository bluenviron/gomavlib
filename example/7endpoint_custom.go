// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

// this is an example struct that implements io.ReadWriteCloser.
// it does not read anything and prints what it receives.
// the only requisite is that Close() must release Read().
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
	if ok == false {
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
	// - understands ardupilotmega dialect
	// - writes messages with given system id and component id
	// - reads/writes to a custom endpoint.
	node, err := gomavlib.NewNode(gomavlib.NodeConf{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointCustom{endpoint},
		},
		Dialect:     ardupilotmega.Dialect,
		SystemId:    10,
		ComponentId: 1,
	})
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// queue a dummy message
	endpoint.readChan <- []byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b")

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
