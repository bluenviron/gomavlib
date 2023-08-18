package gomavlib

import (
	"fmt"
)

type channelAccepter struct {
	n   *Node
	eca endpointChannelAccepter
}

func newChannelAccepter(n *Node, eca endpointChannelAccepter) (*channelAccepter, error) {
	return &channelAccepter{
		n:   n,
		eca: eca,
	}, nil
}

func (ca *channelAccepter) close() {
	ca.eca.close()
}

func (ca *channelAccepter) start() {
	ca.n.channelAcceptersWg.Add(1)
	go ca.run()
}

func (ca *channelAccepter) run() {
	defer ca.n.channelAcceptersWg.Done()

	for {
		label, rwc, err := ca.eca.accept()
		if err != nil {
			if err != errTerminated {
				panic("errTerminated is the only error allowed here")
			}
			break
		}

		ch, err := newChannel(ca.n, ca.eca, label, rwc)
		if err != nil {
			panic(fmt.Errorf("newChannel unexpected error: %s", err))
		}

		ca.n.newChannel(ch)
	}
}
