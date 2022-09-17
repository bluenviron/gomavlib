package gomavlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/message"
)

func TestNodeStreamRequest(t *testing.T) {
	func() {
		node1, err := NewNode(NodeConf{
			Dialect: &dialect.Dialect{3, []message.Message{ //nolint:govet
				&MessageHeartbeat{},
				&MessageRequestDataStream{},
			}},
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

		node2, err := NewNode(NodeConf{
			Dialect: &dialect.Dialect{3, []message.Message{ //nolint:govet
				&MessageHeartbeat{},
				&MessageRequestDataStream{},
			}},
			OutVersion:  V2,
			OutSystemID: 10,
			Endpoints: []EndpointConf{
				EndpointUDPClient{"127.0.0.1:5600"},
			},
			HeartbeatDisable:       false,
			HeartbeatPeriod:        500 * time.Millisecond,
			HeartbeatAutopilotType: 3, // MAV_AUTOPILOT_ARDUPILOTMEGA
		})
		require.NoError(t, err)
		defer node2.Close()

		go func() {
			for range node1.Events() {
			}
		}()

		for evt := range node2.Events() {
			if ee, ok := evt.(*EventFrame); ok {
				if _, ok = ee.Message().(*MessageRequestDataStream); ok {
					return
				}
			}
		}
	}()
}
