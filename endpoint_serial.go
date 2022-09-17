package gomavlib

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/tarm/serial"

	"github.com/aler9/gomavlib/pkg/multibuffer"
)

const (
	serialReconnectPeriod = 2 * time.Second
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

	ctx         context.Context
	ctxCancel   func()
	name        string
	baud        int
	mb          *multibuffer.MultiBuffer
	writerMutex sync.Mutex
	writer      io.Writer

	// in
	read chan []byte
}

func (conf EndpointSerial) init() (Endpoint, error) {
	matches := reSerial.FindStringSubmatch(conf.Address)
	if matches == nil {
		return nil, fmt.Errorf("invalid address")
	}

	name := matches[1]
	baud, _ := strconv.Atoi(matches[2])

	// check device existence
	test, err := serial.OpenPort(&serial.Config{
		Name: name,
		Baud: baud,
	})
	if err != nil {
		return nil, err
	}
	test.Close()

	ctx, ctxCancel := context.WithCancel(context.Background())

	t := &endpointSerial{
		conf:      conf,
		ctx:       ctx,
		ctxCancel: ctxCancel,
		name:      name,
		baud:      baud,
		mb:        multibuffer.New(2, bufferSize),
		read:      make(chan []byte),
	}

	go t.run()

	return t, nil
}

func (t *endpointSerial) isEndpoint() {}

func (t *endpointSerial) Conf() EndpointConf {
	return t.conf
}

func (t *endpointSerial) label() string {
	return "serial"
}

func (t *endpointSerial) Close() error {
	t.ctxCancel()
	return nil
}

func (t *endpointSerial) run() {
	for {
		t.runInner()

		select {
		case <-time.After(serialReconnectPeriod):
		case <-t.ctx.Done():
			return
		}
	}
}

func (t *endpointSerial) runInner() error {
	ser, err := serial.OpenPort(&serial.Config{
		Name: t.name,
		Baud: t.baud,
	})
	if err != nil {
		return err
	}

	func() {
		t.writerMutex.Lock()
		defer t.writerMutex.Unlock()
		t.writer = ser
	}()

	readDone := make(chan error)
	go func() {
		readDone <- func() error {
			for {
				buf := t.mb.Next()
				n, err := ser.Read(buf)
				if err != nil {
					return err
				}

				select {
				case t.read <- buf[:n]:
				case <-t.ctx.Done():
					return errTerminated
				}
			}
		}()
	}()

	select {
	case err := <-readDone:
		ser.Close()
		return err

	case <-t.ctx.Done():
		ser.Close()
		<-readDone
		return errTerminated
	}
}

func (t *endpointSerial) Read(buf []byte) (int, error) {
	select {
	case src := <-t.read:
		n := copy(buf, src)
		return n, nil

	case <-t.ctx.Done():
		return 0, errTerminated
	}
}

func (t *endpointSerial) Write(buf []byte) (int, error) {
	t.writerMutex.Lock()
	defer t.writerMutex.Unlock()

	if t.writer == nil {
		return 0, fmt.Errorf("disconnected")
	}

	return t.writer.Write(buf)
}
