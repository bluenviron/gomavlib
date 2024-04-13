//autogenerated:yes
//nolint:revive,misspell,govet,lll
package matrixpilot

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Flight information.
// This includes time since boot for arm, takeoff, and land, and a flight number.
// Takeoff and landing values reset to zero on arm.
// This can be requested using MAV_CMD_REQUEST_MESSAGE.
// Note, some fields are misnamed - timestamps are from boot (not UTC) and the flight_uuid is a sequence number.
type MessageFlightInformation = common.MessageFlightInformation
