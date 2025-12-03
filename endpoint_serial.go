package gomavlib

import (
	"context"
	"io"
	"net"
	"time"

	"go.bug.st/serial"
)

var serialOpenFunc = func(device string, baud int) (io.ReadWriteCloser, error) {
	dev, err := serial.Open(device, &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		return nil, err
	}

	dev.SetDTR(true) //nolint:errcheck
	dev.SetRTS(true) //nolint:errcheck

	return dev, nil
}

type rwcToConn struct {
	io.ReadWriteCloser
}

func (*rwcToConn) LocalAddr() net.Addr {
	return nil
}

func (*rwcToConn) RemoteAddr() net.Addr {
	return nil
}

func (*rwcToConn) SetDeadline(_ time.Time) error {
	return nil
}

func (*rwcToConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (*rwcToConn) SetWriteDeadline(_ time.Time) error {
	return nil
}

// EndpointSerial sets up a endpoint that works with a serial port.
type EndpointSerial struct {
	// name of the device of the serial port (i.e: /dev/ttyUSB0)
	Device string

	// baud rate (i.e: 57600)
	Baud int
}

func (conf EndpointSerial) init(node *Node) (Endpoint, error) {
	e := &endpointClient{
		node: node,
		conf: EndpointCustomClient{
			Connect: func(_ context.Context) (net.Conn, error) {
				rwc, err := serialOpenFunc(conf.Device, conf.Baud)
				if err != nil {
					return nil, err
				}
				return &rwcToConn{rwc}, nil
			},
			Label: "serial:" + conf.Device,
		},
	}
	err := e.initialize()
	return e, err
}
