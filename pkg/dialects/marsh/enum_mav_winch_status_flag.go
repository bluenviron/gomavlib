//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Winch status flags used in WINCH_STATUS
type MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_FLAG

const (
	// Winch is healthy
	MAV_WINCH_STATUS_HEALTHY MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_HEALTHY
	// Winch line is fully retracted
	MAV_WINCH_STATUS_FULLY_RETRACTED MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_FULLY_RETRACTED
	// Winch motor is moving
	MAV_WINCH_STATUS_MOVING MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_MOVING
	// Winch clutch is engaged allowing motor to move freely.
	MAV_WINCH_STATUS_CLUTCH_ENGAGED MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_CLUTCH_ENGAGED
	// Winch is locked by locking mechanism.
	MAV_WINCH_STATUS_LOCKED MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_LOCKED
	// Winch is gravity dropping payload.
	MAV_WINCH_STATUS_DROPPING MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_DROPPING
	// Winch is arresting payload descent.
	MAV_WINCH_STATUS_ARRESTING MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_ARRESTING
	// Winch is using torque measurements to sense the ground.
	MAV_WINCH_STATUS_GROUND_SENSE MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_GROUND_SENSE
	// Winch is returning to the fully retracted position.
	MAV_WINCH_STATUS_RETRACTING MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_RETRACTING
	// Winch is redelivering the payload. This is a failover state if the line tension goes above a threshold during RETRACTING.
	MAV_WINCH_STATUS_REDELIVER MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_REDELIVER
	// Winch is abandoning the line and possibly payload. Winch unspools the entire calculated line length. This is a failover state from REDELIVER if the number of attempts exceeds a threshold.
	MAV_WINCH_STATUS_ABANDON_LINE MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_ABANDON_LINE
	// Winch is engaging the locking mechanism.
	MAV_WINCH_STATUS_LOCKING MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_LOCKING
	// Winch is spooling on line.
	MAV_WINCH_STATUS_LOAD_LINE MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_LOAD_LINE
	// Winch is loading a payload.
	MAV_WINCH_STATUS_LOAD_PAYLOAD MAV_WINCH_STATUS_FLAG = common.MAV_WINCH_STATUS_LOAD_PAYLOAD
)
