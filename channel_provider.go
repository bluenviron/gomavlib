package gomavlib

import (
	"fmt"
)

type channelProvider struct {
	n   *Node
	eca endpointChannelProvider
}

func newChannelProvider(n *Node, eca endpointChannelProvider) (*channelProvider, error) {
	return &channelProvider{
		n:   n,
		eca: eca,
	}, nil
}

func (ca *channelProvider) close() {
	ca.eca.close()
}

func (ca *channelProvider) start() {
	ca.n.wg.Add(1)
	go ca.run()
}

func (ca *channelProvider) run() {
	defer ca.n.wg.Done()

	for {
		label, rwc, err := ca.eca.provide()
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
