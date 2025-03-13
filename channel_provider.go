package gomavlib

import (
	"errors"
	"fmt"
)

type channelProvider struct {
	node *Node
	eca  endpointChannelProvider

	terminate chan struct{}
}

func (cp *channelProvider) initialize() error {
	cp.terminate = make(chan struct{})
	return nil
}

func (cp *channelProvider) close() {
	close(cp.terminate)
	cp.eca.close()
}

func (cp *channelProvider) start() {
	cp.node.wg.Add(1)
	go cp.run()
}

func (cp *channelProvider) run() {
	defer cp.node.wg.Done()

	for {
		label, rwc, err := cp.eca.provide()
		if err != nil {
			if !errors.Is(err, errTerminated) {
				panic("errTerminated is the only error allowed here")
			}
			break
		}

		ch := &Channel{
			node:     cp.node,
			endpoint: cp.eca,
			label:    label,
			rwc:      rwc,
		}
		err = ch.initialize()
		if err != nil {
			panic(fmt.Errorf("newChannel unexpected error: %w", err))
		}

		cp.node.newChannel(ch)

		if cp.eca.oneChannelAtAtime() {
			// wait the channel to emit EventChannelClose
			// before creating another channel
			select {
			case <-ch.done:
			case <-cp.terminate:
				return
			}
		}
	}
}
