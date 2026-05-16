package gomavlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNodeHeartbeat(t *testing.T) {
	node1 := &Node{
		Dialect:          testDialect,
		OutVersion:       V2,
		OutSystemID:      10,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5600"}},
		HeartbeatDisable: true,
	}
	err := node1.Initialize()
	require.NoError(t, err)
	defer node1.Close()

	node2 := &Node{
		Dialect:         testDialect,
		OutVersion:      V2,
		OutSystemID:     11,
		Endpoints:       []EndpointConf{EndpointUDPClient{"127.0.0.1:5600"}},
		HeartbeatPeriod: 500 * time.Millisecond,
	}
	err = node2.Initialize()
	require.NoError(t, err)
	defer node2.Close()

	<-node1.Events()
	evt := <-node1.Events()
	fr, ok := evt.(*EventFrame)
	require.Equal(t, true, ok)
	_, ok = fr.Message().(*MessageHeartbeat)
	require.Equal(t, true, ok)
}
