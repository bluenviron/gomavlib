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

func (ca *channelAccepter) run() {
	for {
		label, rwc, err := ca.eca.Accept()
		if err != nil {
			if err != errorTerminated {
				panic("errorTerminated is the only error allowed here")
			}
			break
		}

		ch := newChannel(ca.n, ca.eca, label, rwc)
		ca.n.eventsIn <- &eventInChannelNew{ch}
	}
}
