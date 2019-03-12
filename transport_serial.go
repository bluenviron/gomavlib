package gomavlib

import (
	"github.com/tarm/serial"
)

// TransportSerial sends and reads frames through a serial port.
type TransportSerial struct {
	// the name or path of the serial port, ie /dev/ttyAMA0 or COM45
	Name string
	// baud rate, ie 57600
	Baud int
}

func (conf TransportSerial) init(node *Node) (transport, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: conf.Name,
		Baud: conf.Baud,
	})
	if err != nil {
		return nil, err
	}

	ts := TransportCustom{
		ReadWriteCloser: port,
	}
	return ts.init(node)
}
