package gomavlib

import (
	"github.com/tarm/serial"
	"io"
)

// EndpointSerial sets up a endpoint that works through a serial port.
type EndpointSerial struct {
	// the name or path of the serial port, example: /dev/ttyAMA0 or COM45
	Name string
	// baud rate, example: 57600
	Baud int
}

type endpointSerial struct {
	io.ReadWriteCloser
}

func (conf EndpointSerial) init() (endpoint, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: conf.Name,
		Baud: conf.Baud,
	})
	if err != nil {
		return nil, err
	}

	t := &endpointSerial{
		ReadWriteCloser: port,
	}
	return t, nil
}

func (*endpointSerial) isEndpoint() {
}
