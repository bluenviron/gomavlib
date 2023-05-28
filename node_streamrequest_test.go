package gomavlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

func TestNodeStreamRequest(t *testing.T) {
	dialect := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageHeartbeat{},
			&MessageRequestDataStream{},
		},
	}

	node1, err := NewNode(NodeConf{
		Dialect:     dialect,
		OutVersion:  V2,
		OutSystemID: 10,
		Endpoints: []EndpointConf{
			EndpointUDPServer{"127.0.0.1:5600"},
		},
		HeartbeatDisable:    true,
		StreamRequestEnable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	go func() {
		for range node1.Events() { //nolint:revive
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:     dialect,
		OutVersion:  V2,
		OutSystemID: 11,
		Endpoints: []EndpointConf{
			EndpointUDPClient{"127.0.0.1:5600"},
		},
		HeartbeatDisable:       false,
		HeartbeatPeriod:        500 * time.Millisecond,
		HeartbeatAutopilotType: 3, // MAV_AUTOPILOT_ARDUPILOTMEGA
	})
	require.NoError(t, err)
	defer node2.Close()

	for evt := range node2.Events() {
		if ee, ok := evt.(*EventFrame); ok {
			if _, ok = ee.Message().(*MessageRequestDataStream); ok {
				return
			}
		}
	}
}
