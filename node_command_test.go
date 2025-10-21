package gomavlib

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

type (
	MAV_CMD    uint64 //nolint:revive
	MAV_RESULT uint64 //nolint:revive
	MAV_FRAME  uint64 //nolint:revive
)

const (
	MAV_RESULT_ACCEPTED             MAV_RESULT = 0
	MAV_RESULT_TEMPORARILY_REJECTED MAV_RESULT = 1
	MAV_RESULT_DENIED               MAV_RESULT = 2
	MAV_RESULT_UNSUPPORTED          MAV_RESULT = 3
	MAV_RESULT_FAILED               MAV_RESULT = 4
	MAV_RESULT_IN_PROGRESS          MAV_RESULT = 5

	MAV_CMD_COMPONENT_ARM_DISARM MAV_CMD = 400
	MAV_CMD_NAV_WAYPOINT         MAV_CMD = 16

	MAV_FRAME_GLOBAL_RELATIVE_ALT MAV_FRAME = 3
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
				if cmd, ok := frm.Message().(*MessageCommandLong); ok {
					node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}) //nolint:errcheck
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

	// Send command - this blocks until response or timeout
	resp, err := node2.SendCommandLong(&CommandLongRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         uint64(MAV_CMD_COMPONENT_ARM_DISARM),
		Params:          [7]float32{1, 0, 0, 0, 0, 0, 0},
		Options: &CommandOptions{
			Timeout: 2 * time.Second,
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
	require.Less(t, resp.ResponseTime, 2*time.Second)
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
				if cmd, ok := frm.Message().(*MessageCommandInt); ok && !responded {
					node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}) //nolint:errcheck
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

	resp, err := node2.SendCommandInt(&CommandIntRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Frame:           uint64(MAV_FRAME_GLOBAL_RELATIVE_ALT),
		Command:         uint64(MAV_CMD_NAV_WAYPOINT),
		X:               -353621474, // lat * 1e7
		Y:               1491651746, // lon * 1e7
		Z:               100,        // altitude
		Params:          [4]float32{0, 10, 0, 0},
		Options: &CommandOptions{
			Timeout: 2 * time.Second,
		},
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

	start := time.Now()
	resp, err := node2.SendCommandLong(&CommandLongRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         uint64(MAV_CMD_COMPONENT_ARM_DISARM),
		Params:          [7]float32{1, 0, 0, 0, 0, 0, 0},
		Options: &CommandOptions{
			Timeout: 1 * time.Second,
		},
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
				if cmd, ok := frm.Message().(*MessageCommandLong); ok && !responded {
					// Send IN_PROGRESS updates
					for progress := uint8(0); progress <= 100; progress += 50 {
						node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
							Command:      cmd.Command,
							Result:       MAV_RESULT_IN_PROGRESS,
							Progress:     progress,
							TargetSystem: frm.SystemID(),
							TargetComp:   frm.ComponentID(),
						}) //nolint:errcheck
						time.Sleep(50 * time.Millisecond)
					}

					// Send final ACK
					node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}) //nolint:errcheck
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

	progressUpdates := []uint8{}
	resp, err := node2.SendCommandLong(&CommandLongRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         uint64(MAV_CMD_COMPONENT_ARM_DISARM),
		Params:          [7]float32{1, 0, 0, 0, 0, 0, 0},
		Options: &CommandOptions{
			Timeout: 2 * time.Second,
			OnProgress: func(progress uint8) {
				progressUpdates = append(progressUpdates, progress)
			},
		},
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
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
				if cmd, ok := frm.Message().(*MessageCommandLong); ok && !responded {
					node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_DENIED,
						ResultParam2: 42, // Some error code
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}) //nolint:errcheck
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

	resp, err := node2.SendCommandLong(&CommandLongRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         uint64(MAV_CMD_COMPONENT_ARM_DISARM),
		Params:          [7]float32{1, 0, 0, 0, 0, 0, 0},
		Options: &CommandOptions{
			Timeout: 2 * time.Second,
		},
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
	_, err = node2.SendCommandLong(&CommandLongRequest{
		Channel:         channelOpen.Channel,
		TargetSystem:    1,
		TargetComponent: 1,
		Command:         uint64(MAV_CMD_COMPONENT_ARM_DISARM),
		Params:          [7]float32{1, 0, 0, 0, 0, 0, 0},
		Options: &CommandOptions{
			Timeout: 1 * time.Second,
		},
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
				if cmd, ok := frm.Message().(*MessageCommandLong); ok {
					// Small delay to simulate processing
					time.Sleep(100 * time.Millisecond)
					node1.WriteMessageTo(frm.Channel, &MessageCommandAck{
						Command:      cmd.Command,
						Result:       MAV_RESULT_ACCEPTED,
						TargetSystem: frm.SystemID(),
						TargetComp:   frm.ComponentID(),
					}) //nolint:errcheck
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
		go func(cmdID MAV_CMD) {
			resp, err := node2.SendCommandLong(&CommandLongRequest{
				Channel:         channelOpen.Channel,
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         uint64(cmdID),
				Params:          [7]float32{0, 0, 0, 0, 0, 0, 0},
				Options: &CommandOptions{
					Timeout: 2 * time.Second,
				},
			})
			if err != nil {
				errors <- err
			} else {
				results <- resp
			}
		}(MAV_CMD(400 + i))
	}

	// Collect all results
	for i := 0; i < numCommands; i++ {
		select {
		case resp := <-results:
			require.Equal(t, uint64(MAV_RESULT_ACCEPTED), resp.Result)
		case err := <-errors:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Fatalf("timeout waiting for command responses (got %d/%d)", i, numCommands)
		}
	}
}
