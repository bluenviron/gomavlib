// Package main contains an example.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// This example shows how to:
// 1. Create a node that supports the command microservice
// 2. Send a COMMAND_LONG and wait for COMMAND_ACK response
// 3. Handle command results including progress updates

func waitForHeartbeat(node *gomavlib.Node) (*gomavlib.Channel, uint8, uint8) {
	for {
		evt := <-node.Events()

		if evt, ok := evt.(*gomavlib.EventFrame); ok {
			if _, ok = evt.Message().(*common.MessageHeartbeat); ok {
				return evt.Channel, evt.SystemID(), evt.ComponentID()
			}
		}
	}
}

func writeAndWaitCommandLong(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	cmd *common.MessageCommandLong,
	timeout time.Duration,
) error {
	err := node.WriteMessageTo(channel, cmd)
	if err != nil {
		return err
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if ack, ok2 := evt.Message().(*common.MessageCommandAck); ok2 {
					if ack.Command == cmd.Command &&
						evt.SystemID() == cmd.TargetSystem &&
						evt.ComponentID() == cmd.TargetComponent {
						switch {
						case ack.Result == common.MAV_RESULT_IN_PROGRESS:
							log.Printf("command progress: %d%%\n", ack.Progress)

						case ack.Result != common.MAV_RESULT_ACCEPTED:
							return fmt.Errorf("command failed with state %v", ack.Result)

						default:
							return nil
						}
					}
				}
			}

		case <-t.C:
			return fmt.Errorf("command timed out")
		}
	}
}

func sleepWhileListening(node *gomavlib.Node, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-node.Events():
		case <-t.C:
			return
		}
	}
}

func main() {
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyUSB0",
				Baud:   57600,
			},
		},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 255, // Ground station
	}

	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	log.Println("waiting for a heartbeat...")

	heartbeatChan, heartbeatSystemID, heartbeatComponentID := waitForHeartbeat(node)

	log.Printf("received heartbeat from system %d, component %d\n",
		heartbeatSystemID, heartbeatComponentID)

	log.Println("Sending ARM command...")

	err = writeAndWaitCommandLong(node,
		heartbeatChan,
		&common.MessageCommandLong{
			TargetSystem:    heartbeatSystemID,
			TargetComponent: heartbeatComponentID,
			Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
			Param1:          1,
			Param2:          0,
			Param3:          0,
			Param4:          0,
			Param5:          0,
			Param6:          0,
			Param7:          0,
		},
		5*time.Second,
	)
	if err != nil {
		panic(err)
	}

	log.Printf("command succeeded")

	sleepWhileListening(node, 2*time.Second)

	log.Println("Sending DISARM command...")

	err = writeAndWaitCommandLong(node,
		heartbeatChan,
		&common.MessageCommandLong{
			TargetSystem:    heartbeatSystemID,
			TargetComponent: heartbeatComponentID,
			Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
			Param1:          0,
			Param2:          0,
			Param3:          0,
			Param4:          0,
			Param5:          0,
			Param6:          0,
			Param7:          0,
		},
		5*time.Second,
	)
	if err != nil {
		panic(err)
	}

	log.Printf("command succeeded")
}
