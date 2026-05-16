package gomavlib

import (
	"io"
	"net"

	"github.com/bluenviron/gomavlib/v3/pkg/timednetconn"
)

type endpointServer struct {
	node *Node
	conf EndpointCustomServer

	listener net.Listener

	// in
	terminate chan struct{}
}

func (e *endpointServer) initialize() error {
	var err error
	e.listener, err = e.conf.Listen()
	if err != nil {
		return err
	}

	e.terminate = make(chan struct{})

	return nil
}

func (e *endpointServer) isEndpoint() {}

func (e *endpointServer) Conf() EndpointConf {
	return e.conf
}

func (e *endpointServer) close() {
	close(e.terminate)
	e.listener.Close()
}

func (e *endpointServer) oneChannelAtAtime() bool {
	return false
}

func (e *endpointServer) isDatagram() bool {
	return e.conf.IsDatagram
}

func (e *endpointServer) provide() (string, io.ReadWriteCloser, error) {
	nconn, err := e.listener.Accept()
	if err != nil {
		// wait termination, do not report errors
		<-e.terminate
		return "", nil, errTerminated
	}

	label := ""

	if e.conf.Label != "" {
		label += e.conf.Label
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
