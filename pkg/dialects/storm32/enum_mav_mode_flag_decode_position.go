//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/minimal"
)

// These values encode the bit positions of the decode position. These values can be used to read the value of a flag bit by combining the base_mode variable with AND with the flag position value. The result will be either 0 or 1, depending on if the flag is set or not.
type MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION

const (
	// First bit:  10000000
	MAV_MODE_FLAG_DECODE_POSITION_SAFETY MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_SAFETY
	// Second bit: 01000000
	MAV_MODE_FLAG_DECODE_POSITION_MANUAL MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_MANUAL
	// Third bit:  00100000
	MAV_MODE_FLAG_DECODE_POSITION_HIL MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_HIL
	// Fourth bit: 00010000
	MAV_MODE_FLAG_DECODE_POSITION_STABILIZE MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_STABILIZE
	// Fifth bit:  00001000
	MAV_MODE_FLAG_DECODE_POSITION_GUIDED MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_GUIDED
	// Sixth bit:   00000100
	MAV_MODE_FLAG_DECODE_POSITION_AUTO MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_AUTO
	// Seventh bit: 00000010
	MAV_MODE_FLAG_DECODE_POSITION_TEST MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_TEST
	// Eighth bit: 00000001
	MAV_MODE_FLAG_DECODE_POSITION_CUSTOM_MODE MAV_MODE_FLAG_DECODE_POSITION = minimal.MAV_MODE_FLAG_DECODE_POSITION_CUSTOM_MODE
)
