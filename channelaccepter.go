package gomavlib

import (
	"fmt"
)

type channelAccepter struct {
	n   *Node
	eca endpointChannelAccepter

	done chan struct{}
}

func newChannelAccepter(n *Node, eca endpointChannelAccepter) (*channelAccepter, error) {
	return &channelAccepter{
		n:    n,
		eca:  eca,
		done: make(chan struct{}),
	}, nil
}

func (ca *channelAccepter) close() {
	ca.eca.Close()
	<-ca.done
}

func (ca *channelAccepter) run() {
	defer close(ca.done)

	for {
		label, rwc, err := ca.eca.Accept()
		if err != nil {
			if err != errorTerminated {
				panic("errorTerminated is the only error allowed here")
			}
			break
		}

		ch, err := newChannel(ca.n, ca.eca, label, rwc)
		if err != nil {
			panic(fmt.Errorf("newChannel unexpected error: %s", err))
		}

		ca.n.channelNew <- ch
	}
}
