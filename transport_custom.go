package gomavlib

import (
	"io"
)

// TransportCustom reads and writes frames through a custom interface
// that provides the Read(), Write() and Close() functions.
type TransportCustom struct {
	// the struct or interface implementing Read(), Write() and Close()
	ReadWriteCloser io.ReadWriteCloser
}

type transportCustom struct {
	io.ReadWriteCloser
}

func (conf TransportCustom) init() (transport, error) {
	t := &transportCustom{
		conf.ReadWriteCloser,
	}
	return t, nil
}

func (t *transportCustom) isTransport() {
}
