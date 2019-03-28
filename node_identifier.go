package gomavlib

import (
	"fmt"
)

// NodeIdentifier is the unique identifier used to detect other nodes in the network.
type NodeIdentifier struct {
	Channel     *EndpointChannel
	SystemId    byte
	ComponentId byte
}

// string implements fmt.Stringer and returns the node label.
func (i NodeIdentifier) String() string {
	return fmt.Sprintf("chan=%s sid=%d cid=%d", i.Channel, i.SystemId, i.ComponentId)
}
