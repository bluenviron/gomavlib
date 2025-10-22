package gomavlib

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

type (
	MAV_CMD    uint64 //nolint:revive
	MAV_RESULT uint64 //nolint:revive
	MAV_FRAME  uint64 //nolint:revive
)

const (
	MAV_RESULT_ACCEPTED             MAV_RESULT = 0 //nolint:revive
	MAV_RESULT_TEMPORARILY_REJECTED MAV_RESULT = 1 //nolint:revive
	MAV_RESULT_DENIED               MAV_RESULT = 2 //nolint:revive
	MAV_RESULT_UNSUPPORTED          MAV_RESULT = 3 //nolint:revive
	MAV_RESULT_FAILED               MAV_RESULT = 4 //nolint:revive
	MAV_RESULT_IN_PROGRESS          MAV_RESULT = 5 //nolint:revive

	MAV_CMD_COMPONENT_ARM_DISARM MAV_CMD = 400 //nolint:revive
	MAV_CMD_NAV_WAYPOINT         MAV_CMD = 16  //nolint:revive

	MAV_FRAME_GLOBAL_RELATIVE_ALT MAV_FRAME = 3 //nolint:revive
)

type MessageCommandLong struct {
	TargetSystem    uint8
	TargetComponent uint8
	Command         MAV_CMD `mavenum:"uint16"`
	Confirmation    uint8
	Param1          float32
	Param2          float32
	Param3          float32
	Param4          float32
	Param5          float32
	Param6          float32
	Param7          float32
}

func (*MessageCommandLong) GetID() uint32 {
	return 76
}

type MessageCommandInt struct {
	TargetSystem    uint8
	TargetComponent uint8
	Frame           MAV_FRAME `mavenum:"uint8"`
	Command         MAV_CMD   `mavenum:"uint16"`
	Current         uint8
	Autocontinue    uint8
	Param1          float32
	Param2          float32
	Param3          float32
	Param4          float32
	X               int32
	Y               int32
	Z               float32
}

func (*MessageCommandInt) GetID() uint32 {
	return 75
}

type MessageCommandAck struct {
	Command      MAV_CMD    `mavenum:"uint16"`
	Result       MAV_RESULT `mavenum:"uint8"`
	Progress     uint8      `mavext:"true"`
	ResultParam2 int32      `mavext:"true"`
	TargetSystem uint8      `mavext:"true"`
	TargetComp   uint8      `mavext:"true"`
}

func (*MessageCommandAck) GetID() uint32 {
	return 77
}

var commandDialect = &dialect.Dialect{
	Version: 3,
	Messages: []message.Message{
		&MessageHeartbeat{},
		&MessageCommandLong{},
		&MessageCommandInt{},
		&MessageCommandAck{},
	},
}

// TestNodeCommandLongSuccess tests successful COMMAND_LONG execution
func TestNodeCommandLongSuccess(t *testing.T) {
	// Create responder node that will send COMMAND_ACK
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		OutComponentID:   1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5600"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Consume node1 events and respond to commands - exits when channel closes
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk {
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
				}
			}
		}
	}()

	// Create sender node
	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		OutComponentID:   1,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5600"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	// Wait for channel to open
	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Drain node2 events in background
	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	// Send command - this blocks until response or timeout
	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1, // ARM
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
	// 2.005s
	require.Less(t, int64(resp.ResponseTime), int64(2005*time.Millisecond))
}

// TestNodeCommandIntSuccess tests successful COMMAND_INT execution
func TestNodeCommandIntSuccess(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		OutComponentID:   1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5602"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdIntCastOk := frm.Message().(*MessageCommandInt); cmdIntCastOk && !responded {
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		OutComponentID:   1,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5602"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	resp, err := node2.SendCommandInt(&common.MessageCommandInt{
		TargetSystem:    1,
		TargetComponent: 1,
		Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT,
		Command:         common.MAV_CMD_NAV_WAYPOINT,
		X:               -353621474, // lat * 1e7
		Y:               1491651746, // lon * 1e7
		Z:               100,
		Param1:          0,
		Param2:          10,
		Param3:          0,
		Param4:          0,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
}

// TestNodeCommandTimeout tests command timeout handling
func TestNodeCommandTimeout(t *testing.T) {
	// Create node that won't respond
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5603"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Don't respond to commands
	go func() {
		for range node1.Events() { //nolint:revive
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5603"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Drain node2 events in background
	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	start := time.Now()
	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 1 * time.Second,
	})
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, resp)
	// Timeout response has Result = 0 and elapsed time ~= timeout
	require.GreaterOrEqual(t, elapsed, 1*time.Second)
	require.Less(t, elapsed, 2*time.Second)
}

// TestNodeCommandInProgress tests progress callback handling
func TestNodeCommandInProgress(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5604"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Send progress updates then final ACK
	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk && !responded {
					// Send IN_PROGRESS updates
					for progress := uint8(0); progress <= 100; progress += 50 {
						require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
							Command:      cmd.Command,
							Result:       MAV_RESULT_IN_PROGRESS,
							Progress:     progress,
							TargetSystem: frm.SystemID(),
							TargetComp:   frm.ComponentID(),
						}))
						time.Sleep(50 * time.Millisecond)
					}

					// Send final ACK
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5604"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Drain node2 events in background
	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	progressUpdates := []uint8{}
	progressMutex := sync.Mutex{}

	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
		OnProgress: func(progress uint8) {
			progressMutex.Lock()
			progressUpdates = append(progressUpdates, progress)
			progressMutex.Unlock()
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
	progressMutex.Lock()
	defer progressMutex.Unlock()
	require.Greater(t, len(progressUpdates), 0, "should have received progress updates")
}

// TestNodeCommandDenied tests command rejection handling
func TestNodeCommandDenied(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5605"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk && !responded {
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_DENIED,
						ResultParam2: 42, // Some error code
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5605"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_DENIED), resp.Result)
	require.Equal(t, int32(42), resp.ResultParam2)
}

// TestNodeCommandWithoutDialect tests that commands fail gracefully without dialect
func TestNodeCommandWithoutDialect(t *testing.T) {
	// Use client-server setup to get a channel
	node1, err := NewNode(NodeConf{
		Dialect:          nil, // No dialect on responder
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5606"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	go func() {
		for range node1.Events() { //nolint:revive
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          nil, // No dialect on sender either
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5606"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Should fail because no dialect means no command manager
	_, err = node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 1 * time.Second,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "command manager not initialized")
}

// TestNodeCommandMultipleSimultaneous tests multiple commands in flight
func TestNodeCommandMultipleSimultaneous(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5607"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Respond to all commands
	commandCount := 0
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk {
					// Small delay to simulate processing
					time.Sleep(100 * time.Millisecond)
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					commandCount++
					if commandCount >= 5 {
						return // Exit after handling all commands
					}
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5607"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Drain node2 events in background
	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	// Send multiple commands concurrently
	const numCommands = 5
	results := make(chan *CommandResponse, numCommands)
	errors := make(chan error, numCommands)

	for i := 0; i < numCommands; i++ {
		go func() {
			var resp *CommandResponse
			var cmdErr error
			resp, cmdErr = node2.SendCommandLong(&common.MessageCommandLong{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
				Param1:          1,
				Param2:          0,
				Param3:          0,
				Param4:          0,
				Param5:          0,
				Param6:          0,
				Param7:          0,
			}, &CommandOptions{
				Channel: channelOpen.Channel,
				Timeout: 2 * time.Second,
			})
			if cmdErr != nil {
				errors <- cmdErr
			} else {
				results <- resp
			}
		}()
	}

	// Collect all results
	for i := 0; i < numCommands; i++ {
		select {
		case resp := <-results:
			require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
		case err = <-errors:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatalf("timeout waiting for command responses (got %d/%d)", i, numCommands)
		}
	}
}

// TestNodeCommandNilOptions tests that nil options returns error
func TestNodeCommandNilOptions(t *testing.T) {
	node, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5608"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	go func() {
		for range node.Events() { //nolint:revive
		}
	}()

	_, err = node.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
	}, nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "options is nil")
}

// TestNodeCommandNilChannel tests that nil channel returns error
func TestNodeCommandNilChannel(t *testing.T) {
	node, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5609"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	go func() {
		for range node.Events() { //nolint:revive
		}
	}()

	_, err = node.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
	}, &CommandOptions{
		Channel: nil,
		Timeout: 1 * time.Second,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "need channel")
}

// TestNodeCommandIntAllResultTypes tests different COMMAND_INT result types
func TestNodeCommandIntAllResultTypes(t *testing.T) {
	testCases := []struct {
		name           string
		result         MAV_RESULT
		port           string
		expectError    bool
		expectedResult uint64
	}{
		{"Accepted", MAV_RESULT_ACCEPTED, "127.0.0.1:5611", false, uint64(MAV_RESULT_ACCEPTED)},
		{"Rejected", MAV_RESULT_TEMPORARILY_REJECTED, "127.0.0.1:5612", false, uint64(MAV_RESULT_TEMPORARILY_REJECTED)},
		{"Unsupported", MAV_RESULT_UNSUPPORTED, "127.0.0.1:5613", false, uint64(MAV_RESULT_UNSUPPORTED)},
		{"Failed", MAV_RESULT_FAILED, "127.0.0.1:5614", false, uint64(MAV_RESULT_FAILED)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			node1, err := NewNode(NodeConf{
				Dialect:          commandDialect,
				OutVersion:       V2,
				OutSystemID:      1,
				Endpoints:        []EndpointConf{EndpointUDPServer{tc.port}},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer node1.Close()

			responded := false
			go func() {
				for evt := range node1.Events() {
					if frm, ok := evt.(*EventFrame); ok {
						if cmd, cmdIntCastOk := frm.Message().(*MessageCommandInt); cmdIntCastOk && !responded {
							require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
								Command:      cmd.Command,
								Result:       tc.result,
								TargetSystem: frm.SystemID(),
								TargetComp:   frm.ComponentID(),
							}))
							responded = true
						}
					}
				}
			}()

			node2, err := NewNode(NodeConf{
				Dialect:          commandDialect,
				OutVersion:       V2,
				OutSystemID:      255,
				Endpoints:        []EndpointConf{EndpointUDPClient{tc.port}},
				HeartbeatDisable: true,
			})
			require.NoError(t, err)
			defer node2.Close()

			evt := <-node2.Events()
			channelOpen, ok := evt.(*EventChannelOpen)
			require.True(t, ok)

			go func() {
				for range node2.Events() { //nolint:revive
				}
			}()

			resp, err := node2.SendCommandInt(&common.MessageCommandInt{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         common.MAV_CMD_NAV_WAYPOINT,
			}, &CommandOptions{
				Channel: channelOpen.Channel,
				Timeout: 1 * time.Second,
			})

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, tc.expectedResult, resp.Result)
			}
		})
	}
}

// TestNodeCommandProgressChannelFull tests progress updates when channel is full
func TestNodeCommandProgressChannelFull(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5616"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Send many IN_PROGRESS updates rapidly
	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk && !responded {
					// Send 100 IN_PROGRESS updates rapidly to try to fill buffer
					for progress := uint8(0); progress < 100; progress++ {
						require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
							Command:      cmd.Command,
							Result:       MAV_RESULT_IN_PROGRESS,
							Progress:     progress,
							TargetSystem: frm.SystemID(),
							TargetComp:   frm.ComponentID(),
						}))
					}

					// Send final ACK
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5616"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	progressCount := 0
	progressMutex := sync.Mutex{}

	// Slow progress handler to cause buffer to fill
	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 3 * time.Second,
		OnProgress: func(progress uint8) {
			progressMutex.Lock()
			progressCount++
			progressMutex.Unlock()
			time.Sleep(50 * time.Millisecond) // Slow handler
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
	progressMutex.Lock()
	defer progressMutex.Unlock()
	// Should have received some progress updates, but not necessarily all due to buffer
	require.Greater(t, progressCount, 0)
}

// TestNodeCommandDefaultTimeout tests that default timeout is applied when not specified
func TestNodeCommandDefaultTimeout(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5617"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	// Don't respond to commands
	go func() {
		for range node1.Events() { //nolint:revive
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5617"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	start := time.Now()
	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 0, // Should use default 5 second timeout
	})
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.NotNil(t, resp)
	// Should timeout with default 5 seconds
	require.GreaterOrEqual(t, elapsed, 5*time.Second)
	require.Less(t, elapsed, 6*time.Second)
}

// TestNodeCommandWithoutProgressCallback tests command without progress callback
func TestNodeCommandWithoutProgressCallback(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5618"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk && !responded {
					// Send IN_PROGRESS updates even though client has no callback
					for progress := uint8(0); progress <= 50; progress += 25 {
						require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
							Command:      cmd.Command,
							Result:       MAV_RESULT_IN_PROGRESS,
							Progress:     progress,
							TargetSystem: frm.SystemID(),
							TargetComp:   frm.ComponentID(),
						}))
						time.Sleep(10 * time.Millisecond)
					}

					// Send final ACK
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5618"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	// No OnProgress callback
	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
}

// TestNodeCommandIntNotAvailable tests SendCommandInt when COMMAND_INT not in dialect
func TestNodeCommandIntNotAvailable(t *testing.T) {
	// Create dialect without COMMAND_INT
	dialectNoInt := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageHeartbeat{},
			&MessageCommandLong{},
			&MessageCommandAck{},
			// No MessageCommandInt
		},
	}

	node, err := NewNode(NodeConf{
		Dialect:          dialectNoInt,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5619"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	go func() {
		for range node.Events() { //nolint:revive
		}
	}()

	// Should fail immediately without needing a channel
	// When COMMAND_INT is missing, nodeCommand doesn't initialize at all
	_, err = node.SendCommandInt(&common.MessageCommandInt{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_NAV_WAYPOINT,
	}, &CommandOptions{
		Channel: &Channel{}, // Dummy channel - won't be used
		Timeout: 1 * time.Second,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "command manager not initialized")
}

// TestNodeCommandLongNotAvailable tests SendCommandLong when COMMAND_LONG not in dialect
func TestNodeCommandLongNotAvailable(t *testing.T) {
	// Create dialect without COMMAND_LONG
	dialectNoLong := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageHeartbeat{},
			&MessageCommandInt{},
			&MessageCommandAck{},
			// No MessageCommandLong
		},
	}

	node, err := NewNode(NodeConf{
		Dialect:          dialectNoLong,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5620"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node.Close()

	go func() {
		for range node.Events() { //nolint:revive
		}
	}()

	// Should fail immediately without needing a channel
	// When COMMAND_LONG is missing, nodeCommand doesn't initialize at all
	_, err = node.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
	}, &CommandOptions{
		Channel: &Channel{}, // Dummy channel - won't be used
		Timeout: 1 * time.Second,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "command manager not initialized")
}

// TestNodeCommandAckWithExtensionFields tests that extension fields are properly handled
func TestNodeCommandAckWithExtensionFields(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5621"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node1.Close()

	responded := false
	go func() {
		for evt := range node1.Events() {
			if frm, ok := evt.(*EventFrame); ok {
				if cmd, cmdLongCastOk := frm.Message().(*MessageCommandLong); cmdLongCastOk && !responded {
					// Send ACK with all extension fields populated
					require.NoError(t, node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						Progress:     75,
						ResultParam2: 12345,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}))
					responded = true
				}
			}
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5621"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)
	defer node2.Close()

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	go func() {
		for range node2.Events() { //nolint:revive
		}
	}()

	resp, err := node2.SendCommandLong(&common.MessageCommandLong{
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
		Param1:          1,
	}, &CommandOptions{
		Channel: channelOpen.Channel,
		Timeout: 2 * time.Second,
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
	require.Equal(t, uint8(75), resp.Progress)
	require.Equal(t, int32(12345), resp.ResultParam2)
}

// TestNodeCommandCancelOnClose tests that cancelAllPending is called when node closes
func TestNodeCommandCancelOnClose(t *testing.T) {
	node1, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      1,
		Endpoints:        []EndpointConf{EndpointUDPServer{"127.0.0.1:5622"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	// Don't respond to commands - they'll be pending when we close
	go func() {
		for range node1.Events() { //nolint:revive
		}
	}()

	node2, err := NewNode(NodeConf{
		Dialect:          commandDialect,
		OutVersion:       V2,
		OutSystemID:      255,
		Endpoints:        []EndpointConf{EndpointUDPClient{"127.0.0.1:5622"}},
		HeartbeatDisable: true,
	})
	require.NoError(t, err)

	evt := <-node2.Events()
	channelOpen, ok := evt.(*EventChannelOpen)
	require.True(t, ok)

	// Send multiple commands that won't get responses
	// These will timeout eventually (testing cancelAllPending coverage)
	numCommands := 3
	for i := 0; i < numCommands; i++ {
		go func() {
			// Use short timeout so test doesn't hang
			_, _ = node2.SendCommandLong(&common.MessageCommandLong{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
				Param1:          1,
			}, &CommandOptions{
				Channel: channelOpen.Channel,
				Timeout: 500 * time.Millisecond,
			})
		}()
	}

	// Give commands time to be sent and become pending
	time.Sleep(50 * time.Millisecond)

	// Close nodes - this calls cancelAllPending to clean up pending commands
	// The test passes if Close() returns (doesn't hang)
	node2.Close()
	node1.Close()
}
