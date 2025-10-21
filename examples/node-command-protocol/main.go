// Package main contains an example.
package main

import (
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// This example shows how to:
// 1. Create a node that supports the command protocol
// 2. Send a COMMAND_LONG and wait for COMMAND_ACK response
// 3. Handle command results including progress updates

func main() {
	// Create a node with command protocol support
	// (automatically enabled when a dialect is provided)
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     common.Dialect, // Required for command protocol
		OutVersion:  gomavlib.V2,
		OutSystemID: 255, // Ground station
	}

	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	var targetChannel *gomavlib.Channel

	log.Println("Waiting for connection...")

	// Wait for a channel to open and send a command
	for evt := range node.Events() {
		switch e := evt.(type) {
		case *gomavlib.EventChannelOpen:
			targetChannel = e.Channel
			log.Printf("Channel opened: %s\n", targetChannel)

		case *gomavlib.EventFrame:
			// Wait for heartbeat to know target system ID
			if _, ok := e.Message().(*common.MessageHeartbeat); ok {
				log.Printf("Received heartbeat from system %d, component %d\n",
					e.SystemID(), e.ComponentID())

				// Send ARM command with progress tracking
				log.Println("Sending ARM command...")
				var resp *gomavlib.CommandResponse
				var cmdErr error
				resp, cmdErr = node.SendCommandLong(&common.MessageCommandLong{
					TargetSystem:    e.SystemID(),
					TargetComponent: e.ComponentID(),
					Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
					Param1:          1,
					Param2:          0,
					Param3:          0,
					Param4:          0,
					Param5:          0,
					Param6:          0,
					Param7:          0,
				}, &gomavlib.CommandOptions{
					Timeout: 5 * time.Second,
					OnProgress: func(progress uint8) {
						if progress == 255 {
							log.Println("Command in progress (progress unknown)")
						} else {
							log.Printf("Command progress: %d%%\n", progress)
						}
					},
				})

				if cmdErr != nil {
					log.Printf("Command failed: %v\n", err)
					return
				}

				// Check result
				log.Printf("Command completed in %v\n", resp.ResponseTime)
				switch resp.Result {
				case uint64(common.MAV_RESULT_ACCEPTED):
					log.Println("✓ Command ACCEPTED - Vehicle armed!")

					// Example: Send another command (DISARM)
					time.Sleep(2 * time.Second)
					log.Println("Sending DISARM command...")
					resp, cmdErr = node.SendCommandLong(&common.MessageCommandLong{
						TargetSystem:    e.SystemID(),
						TargetComponent: e.ComponentID(),
						Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
						Param1:          0,
						Param2:          0,
						Param3:          0,
						Param4:          0,
						Param5:          0,
						Param6:          0,
						Param7:          0,
					}, &gomavlib.CommandOptions{
						Timeout: 5 * time.Second,
					})

					if cmdErr != nil {
						log.Printf("DISARM command failed: %v\n", err)
						return
					}

					if resp.Result == uint64(common.MAV_RESULT_ACCEPTED) {
						log.Println("✓ Vehicle disarmed!")
					} else {
						log.Printf("DISARM command failed: %v\n", resp.Result)
					}

					return

				case uint64(common.MAV_RESULT_TEMPORARILY_REJECTED):
					log.Println("✗ Command TEMPORARILY REJECTED - retry later")
					log.Printf("  Additional info: %d\n", resp.ResultParam2)

				case uint64(common.MAV_RESULT_DENIED):
					log.Println("✗ Command DENIED")
					log.Printf("  Additional info: %d\n", resp.ResultParam2)

				case uint64(common.MAV_RESULT_UNSUPPORTED):
					log.Println("✗ Command UNSUPPORTED")

				case uint64(common.MAV_RESULT_FAILED):
					log.Println("✗ Command FAILED")
					log.Printf("  Additional info: %d\n", resp.ResultParam2)

				default:
					log.Printf("✗ Unknown result: %d\n", resp.Result)
				}

				return
			}
		}
	}
}
