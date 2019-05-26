package gomavlib

type channelAccepter struct {
	n   *Node
	eca endpointChannelAccepter
}

func newChannelAccepter(n *Node, eca endpointChannelAccepter) *channelAccepter {
	ca := &channelAccepter{
		n:   n,
		eca: eca,
	}
	return ca
}

func (ca *channelAccepter) close() {
	ca.eca.Close()
}

func (ca *channelAccepter) start() {
	ca.n.wg.Add(1)

	go func() {
		defer ca.n.wg.Done()

		for {
			label, rwc, err := ca.eca.Accept()
			if err != nil {
				if err != errorTerminated {
					panic("errorTerminated is the only error allowed here")
				}
				break
			}

			ch := ca.n.createChannel(ca.eca, label, rwc)
			ch.start()
		}
	}()
}
