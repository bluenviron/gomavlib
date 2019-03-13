package gomavlib

import (
	"github.com/tarm/serial"
	"io"
)

// TransportSerial reads and writes frames through a serial port.
type TransportSerial struct {
	// the name or path of the serial port, example: /dev/ttyAMA0 or COM45
	Name string
	// baud rate, example: 57600
	Baud int
}

type transportSerial struct {
	io.ReadWriteCloser
}

func (conf TransportSerial) init() (transport, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: conf.Name,
		Baud: conf.Baud,
	})
	if err != nil {
		return nil, err
	}

	t := &transportSerial{
		ReadWriteCloser: port,
	}
	return t, nil
}

func (*transportSerial) isTransport() {
}
