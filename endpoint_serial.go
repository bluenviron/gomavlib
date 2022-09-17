package gomavlib

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/tarm/serial"

	"github.com/aler9/gomavlib/pkg/multibuffer"
)

const (
	serialReconnectPeriod = 2 * time.Second
)

// EndpointSerial sets up a endpoint that works with a serial port.
type EndpointSerial struct {
	// the name of the device of the serial port (i.e: /dev/ttyUSB0)
	Device string

	// baud rate (i.e: 57600)
	Baud int
}

type endpointSerial struct {
	conf EndpointSerial

	ctx         context.Context
	ctxCancel   func()
	mb          *multibuffer.MultiBuffer
	writerMutex sync.Mutex
	writer      io.Writer

	// in
	read chan []byte
}

func (conf EndpointSerial) init() (Endpoint, error) {
	// check device existence
	test, err := serial.OpenPort(&serial.Config{
		Name: conf.Device,
		Baud: conf.Baud,
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
		Name: t.conf.Device,
		Baud: t.conf.Baud,
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
