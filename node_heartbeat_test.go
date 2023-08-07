package gomavlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNodeHeartbeat(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	node2, err := NewNode(NodeConf{
		Dialect:     testDialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable: false,
		HeartbeatPeriod:  500 * time.Millisecond,
	})
	require.NoError(t, err)
	defer node2.Close()

	<-node1.Events()
	evt := <-node1.Events()
	fr, ok := evt.(*EventFrame)
	require.Equal(t, true, ok)
	_, ok = fr.Message().(*MessageHeartbeat)
	require.Equal(t, true, ok)
}
