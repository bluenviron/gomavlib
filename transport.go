package gomavlib

import (
	"io"
	"time"
)

const (
	// constant for ip-based transports
	netBufferSize      = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
	netConnectTimeout  = 10 * time.Second
	netReconnectPeriod = 2 * time.Second
	netReadTimeout     = 60 * time.Second
	netWriteTimeout    = 10 * time.Second
)

// TransportChannel is a channel provided by a transport.
type TransportChannel struct {
	rwc       io.ReadWriteCloser
	writeChan chan interface{}
}

// TransportConf is the interface implemented by all transports.
type TransportConf interface {
	init() (transport, error)
}

// a transport must also implement one of the following:
// - transportChannelSingle
// - transportChannelAccepter
type transport interface {
	isTransport()
}

// transport that provides a single channel.
// Read() must not return any error unless Close() is called.
type transportChannelSingle interface {
	transport
	io.ReadWriteCloser
}

// transport that provides multiple channels.
type transportChannelAccepter interface {
	transport
	Close() error
	Accept() (io.ReadWriteCloser, error)
}
