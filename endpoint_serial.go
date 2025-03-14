package gomavlib

import (
	"context"
	"io"
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

	ctx       context.Context
	ctxCancel func()
}

func (e *endpointSerial) initialize() error {
	// check device existence
	test, err := serialOpenFunc(e.conf.Device, e.conf.Baud)
	if err != nil {
		return err
	}
	test.Close()

	e.ctx, e.ctxCancel = context.WithCancel(context.Background())

	return nil
}

func (e *endpointSerial) isEndpoint() {}

func (e *endpointSerial) Conf() EndpointConf {
	return e.conf
}

func (e *endpointSerial) close() {
	e.ctxCancel()
}

func (e *endpointSerial) oneChannelAtAtime() bool {
	return true
}

func (e *endpointSerial) connect() (io.ReadWriteCloser, error) {
	return serialOpenFunc(e.conf.Device, e.conf.Baud)
}

func (e *endpointSerial) provide() (string, io.ReadWriteCloser, error) {
	for {
		conn, err := e.connect()
		if err != nil {
			select {
			case <-time.After(reconnectPeriod):
				continue
			case <-e.ctx.Done():
				return "", nil, errTerminated
			}
		}

		return "serial", conn, nil
	}
}
