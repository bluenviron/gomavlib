package gomavlib

import (
	"context"
	"io"

	"go.bug.st/serial"

	"github.com/chrisdalke/gomavlib/v3/pkg/reconnector"
)

var serialOpenFunc = func(device string, baud int) (io.ReadWriteCloser, error) {
	return serial.Open(device, &serial.Mode{
		BaudRate: baud,
	})
}

// EndpointSerial sets up a endpoint that works with a serial port.
type EndpointSerial struct {
	// the name of the device of the serial port (i.e: /dev/ttyUSB0)
	Device string

	// baud rate (i.e: 57600)
	Baud int
}

type endpointSerial struct {
	conf        EndpointConf
	reconnector *reconnector.Reconnector
}

func (conf EndpointSerial) init(_ *Node) (Endpoint, error) {
	// check device existence
	test, err := serialOpenFunc(conf.Device, conf.Baud)
	if err != nil {
		return nil, err
	}
	test.Close()

	t := &endpointSerial{
		conf: conf,
		reconnector: reconnector.New(
			func(_ context.Context) (io.ReadWriteCloser, error) {
				return serialOpenFunc(conf.Device, conf.Baud)
			},
		),
	}

	return t, nil
}

func (t *endpointSerial) isEndpoint() {}

func (t *endpointSerial) Conf() EndpointConf {
	return t.conf
}

func (t *endpointSerial) close() {
	t.reconnector.Close()
}

func (t *endpointSerial) oneChannelAtAtime() bool {
	return true
}

func (t *endpointSerial) provide() (string, io.ReadWriteCloser, error) {
	conn, ok := t.reconnector.Reconnect()
	if !ok {
		return "", nil, errTerminated
	}

	return "serial", conn, nil
}
