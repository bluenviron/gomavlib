//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Result of mission operation (in a MISSION_ACK message).
type MAV_MISSION_RESULT = common.MAV_MISSION_RESULT

const (
	// mission accepted OK
	MAV_MISSION_ACCEPTED MAV_MISSION_RESULT = common.MAV_MISSION_ACCEPTED
	// Generic error / not accepting mission commands at all right now.
	MAV_MISSION_ERROR MAV_MISSION_RESULT = common.MAV_MISSION_ERROR
	// Coordinate frame is not supported.
	MAV_MISSION_UNSUPPORTED_FRAME MAV_MISSION_RESULT = common.MAV_MISSION_UNSUPPORTED_FRAME
	// Command is not supported.
	MAV_MISSION_UNSUPPORTED MAV_MISSION_RESULT = common.MAV_MISSION_UNSUPPORTED
	// Mission items exceed storage space.
	MAV_MISSION_NO_SPACE MAV_MISSION_RESULT = common.MAV_MISSION_NO_SPACE
	// One of the parameters has an invalid value.
	MAV_MISSION_INVALID MAV_MISSION_RESULT = common.MAV_MISSION_INVALID
	// param1 has an invalid value.
	MAV_MISSION_INVALID_PARAM1 MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM1
	// param2 has an invalid value.
	MAV_MISSION_INVALID_PARAM2 MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM2
	// param3 has an invalid value.
	MAV_MISSION_INVALID_PARAM3 MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM3
	// param4 has an invalid value.
	MAV_MISSION_INVALID_PARAM4 MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM4
	// x / param5 has an invalid value.
	MAV_MISSION_INVALID_PARAM5_X MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM5_X
	// y / param6 has an invalid value.
	MAV_MISSION_INVALID_PARAM6_Y MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM6_Y
	// z / param7 has an invalid value.
	MAV_MISSION_INVALID_PARAM7 MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_PARAM7
	// Mission item received out of sequence
	MAV_MISSION_INVALID_SEQUENCE MAV_MISSION_RESULT = common.MAV_MISSION_INVALID_SEQUENCE
	// Not accepting any mission commands from this communication partner.
	MAV_MISSION_DENIED MAV_MISSION_RESULT = common.MAV_MISSION_DENIED
	// Current mission operation cancelled (e.g. mission upload, mission download).
	MAV_MISSION_OPERATION_CANCELLED MAV_MISSION_RESULT = common.MAV_MISSION_OPERATION_CANCELLED
)
