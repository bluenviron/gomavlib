package gomavlib

import (
	"context"
	"io"

	"go.bug.st/serial"

	"github.com/bluenviron/gomavlib/v3/pkg/reconnector"
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

// EndpointSerial sets up a endpoint that works with a serial port.
type EndpointSerial struct {
	// name of the device of the serial port (i.e: /dev/ttyUSB0)
	Device string

	// baud rate (i.e: 57600)
	Baud int
}

func (conf EndpointSerial) init(node *Node) (Endpoint, error) {
	e := &endpointSerial{
		node: node,
		conf: conf,
	}
	err := e.initialize()
	return e, err
}

type endpointSerial struct {
	node *Node
	conf EndpointSerial

	reconnector *reconnector.Reconnector
}

func (e *endpointSerial) initialize() error {
	// check device existence
	test, err := serialOpenFunc(e.conf.Device, e.conf.Baud)
	if err != nil {
		return err
	}
	test.Close()

	e.reconnector = reconnector.New(
		func(_ context.Context) (io.ReadWriteCloser, error) {
			return serialOpenFunc(e.conf.Device, e.conf.Baud)
		},
	)

	return nil
}

func (e *endpointSerial) isEndpoint() {}

func (e *endpointSerial) Conf() EndpointConf {
	return e.conf
}

func (e *endpointSerial) close() {
	e.reconnector.Close()
}

func (e *endpointSerial) oneChannelAtAtime() bool {
	return true
}

func (e *endpointSerial) provide() (string, io.ReadWriteCloser, error) {
	conn, ok := e.reconnector.Reconnect()
	if !ok {
		return "", nil, errTerminated
	}

	return "serial", conn, nil
}
