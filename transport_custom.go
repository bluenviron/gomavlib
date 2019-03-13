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
	conf      TransportCustom
	node      *Node
	terminate chan struct{}
	tconn     *TransportChannel
}

func (conf TransportCustom) init(node *Node) (transport, error) {
	t := &transportCustom{
		conf:      conf,
		node:      node,
		terminate: make(chan struct{}),
	}

	t.tconn = &TransportChannel{
		transport: t,
		writer:    t.conf.ReadWriteCloser,
	}
	t.node.addTransportChannel(t.tconn)

	return t, nil
}

func (t *transportCustom) closePrematurely() {
	t.conf.ReadWriteCloser.Close()
}

func (t *transportCustom) do() {
	listenDone := make(chan struct{})
	go func() {
		defer func() { listenDone <- struct{}{} }()
		defer t.node.removeTransportChannel(t.tconn)

		var buf [netBufferSize]byte
		for {
			n, err := t.conf.ReadWriteCloser.Read(buf[:])
			if err != nil {
				break
			}
			t.node.processBuffer(nil, buf[:n])
		}
	}()

	select {
	// unexpected error, wait for terminate
	case <-listenDone:
		t.conf.ReadWriteCloser.Close()
		<-t.node.terminate

	// terminated
	case <-t.node.terminate:
		t.conf.ReadWriteCloser.Close()
		<-listenDone
	}
}
