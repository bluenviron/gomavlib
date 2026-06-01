package gomavlib

import (
	"io"
	"net"

	"github.com/bluenviron/gomavlib/v4/pkg/timednetconn"
)

var _ Endpoint = (*EndpointCustomServer)(nil)

// EndpointCustomServer is an endpoint that works with custom implementations
// by providing a custom Listen func that returns a net.Listener.
// This allows you to use custom protocols that conform to the net.listner.
// A use case could be to add encrypted protocol implementations like DTLS or TCP with TLS.
type EndpointCustomServer struct {
	// function to invoke when server should start listening
	Listen func() (net.Listener, error)

	// the label of the protocol
	Label string

	// whether the connection is datagram-based (e.g. UDP).
	IsDatagram bool

	node      *Node
	listener  net.Listener
	terminate chan struct{}
}

func (e *EndpointCustomServer) init(node *Node) error {
	e.node = node

	var err error
	e.listener, err = e.Listen()
	if err != nil {
		return err
	}

	e.terminate = make(chan struct{})

	return err
}

func (e *EndpointCustomServer) isEndpoint() {}

func (e *EndpointCustomServer) close() {
	close(e.terminate)
	e.listener.Close()
}

func (e *EndpointCustomServer) oneChannelAtAtime() bool {
	return false
}

func (e *EndpointCustomServer) isDatagram() bool {
	return e.IsDatagram
}

func (e *EndpointCustomServer) provide() (string, io.ReadWriteCloser, error) {
	nconn, err := e.listener.Accept()
	if err != nil {
		// wait termination, do not report errors
		<-e.terminate
		return "", nil, errTerminated
	}

	label := ""

	if e.Label != "" {
		label += e.Label
	} else {
		label = "custom"
	}

	label += ":" + nconn.RemoteAddr().String()

	conn := timednetconn.New(
		e.node.IdleTimeout,
		e.node.WriteTimeout,
		nconn)

	return label, conn, nil
}
