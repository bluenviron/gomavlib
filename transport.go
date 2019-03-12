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

type TransportConf interface {
	init(n *Node) (transport, error)
}

type transport interface {
	closePrematurely()
	do()
}

type TransportChannel struct {
	transport transport
	writer    io.Writer
	writeChan chan interface{}
}
