//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Bitmap to indicate which dimensions should be ignored by the vehicle: a value of 0b0000000000000000 or 0b0000001000000000 indicates that none of the setpoint dimensions should be ignored. If bit 9 is set the floats afx afy afz should be interpreted as force instead of acceleration.
type POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK

const (
	// Ignore position x
	POSITION_TARGET_TYPEMASK_X_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_X_IGNORE
	// Ignore position y
	POSITION_TARGET_TYPEMASK_Y_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_Y_IGNORE
	// Ignore position z
	POSITION_TARGET_TYPEMASK_Z_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_Z_IGNORE
	// Ignore velocity x
	POSITION_TARGET_TYPEMASK_VX_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_VX_IGNORE
	// Ignore velocity y
	POSITION_TARGET_TYPEMASK_VY_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_VY_IGNORE
	// Ignore velocity z
	POSITION_TARGET_TYPEMASK_VZ_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_VZ_IGNORE
	// Ignore acceleration x
	POSITION_TARGET_TYPEMASK_AX_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_AX_IGNORE
	// Ignore acceleration y
	POSITION_TARGET_TYPEMASK_AY_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_AY_IGNORE
	// Ignore acceleration z
	POSITION_TARGET_TYPEMASK_AZ_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_AZ_IGNORE
	// Use force instead of acceleration
	POSITION_TARGET_TYPEMASK_FORCE_SET POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_FORCE_SET
	// Ignore yaw
	POSITION_TARGET_TYPEMASK_YAW_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_YAW_IGNORE
	// Ignore yaw rate
	POSITION_TARGET_TYPEMASK_YAW_RATE_IGNORE POSITION_TARGET_TYPEMASK = common.POSITION_TARGET_TYPEMASK_YAW_RATE_IGNORE
)
