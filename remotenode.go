package gomavlib

import (
	"fmt"
)

// RemoteNode is the unique identifier used to detect other nodes in the network.
type RemoteNode struct {
	Channel     *Channel
	SystemId    byte
	ComponentId byte
}

// string implements fmt.Stringer and returns the node label.
func (i RemoteNode) String() string {
	return fmt.Sprintf("chan=%s sid=%d cid=%d", i.Channel, i.SystemId, i.ComponentId)
}
