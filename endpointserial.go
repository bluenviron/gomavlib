package gomavlib

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/tarm/serial"
)

var reSerial = regexp.MustCompile("^(.+?):([0-9]+)$")

// EndpointSerial sets up a endpoint that works with a serial port.
type EndpointSerial struct {
	// the address of the serial port in format name:baudrate
	// example: /dev/ttyUSB0:57600
	Address string
}

type endpointSerial struct {
	conf EndpointSerial
	io.ReadWriteCloser
}

func (conf EndpointSerial) init() (Endpoint, error) {
	matches := reSerial.FindStringSubmatch(conf.Address)
	if matches == nil {
		return nil, fmt.Errorf("invalid address")
	}

	name := matches[1]
	baud, _ := strconv.Atoi(matches[2])

	rwc, err := serial.OpenPort(&serial.Config{
		Name: name,
		Baud: baud,
	})
	if err != nil {
		return nil, err
	}

	t := &endpointSerial{
		conf:            conf,
		ReadWriteCloser: rwc,
	}
	return t, nil
}

func (t *endpointSerial) isEndpoint() {}

func (t *endpointSerial) Conf() EndpointConf {
	return t.conf
}

func (t *endpointSerial) Label() string {
	return "serial"
}
