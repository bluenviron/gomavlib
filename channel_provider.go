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

func (cp *channelProvider) close() {
	cp.eca.close()
}

func (cp *channelProvider) start() {
	cp.n.wg.Add(1)
	go cp.run()
}

func (cp *channelProvider) run() {
	defer cp.n.wg.Done()

	for {
		label, rwc, err := cp.eca.provide()
		if err != nil {
			if err != errTerminated {
				panic("errTerminated is the only error allowed here")
			}
			break
		}

		ch, err := newChannel(cp.n, cp.eca, label, rwc)
		if err != nil {
			panic(fmt.Errorf("newChannel unexpected error: %s", err))
		}

		cp.n.newChannel(ch)

		if cp.eca.oneChannelAtAtime() {
			// wait the channel to emit EventChannelClose
			// before creating another channel
			select {
			case <-ch.done:
			case <-cp.n.terminate:
			}
		}
	}
}
