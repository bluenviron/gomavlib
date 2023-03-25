//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package cubepilot

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Focus types for MAV_CMD_SET_CAMERA_FOCUS
type SET_FOCUS_TYPE = common.SET_FOCUS_TYPE

const (
	// Focus one step increment (-1 for focusing in, 1 for focusing out towards infinity).
	FOCUS_TYPE_STEP SET_FOCUS_TYPE = common.FOCUS_TYPE_STEP
	// Continuous focus up/down until stopped (-1 for focusing in, 1 for focusing out towards infinity, 0 to stop focusing)
	FOCUS_TYPE_CONTINUOUS SET_FOCUS_TYPE = common.FOCUS_TYPE_CONTINUOUS
	// Focus value as proportion of full camera focus range (a value between 0.0 and 100.0)
	FOCUS_TYPE_RANGE SET_FOCUS_TYPE = common.FOCUS_TYPE_RANGE
	// Focus value in metres. Note that there is no message to get the valid focus range of the camera, so this can type can only be used for cameras where the range is known (implying that this cannot reliably be used in a GCS for an arbitrary camera).
	FOCUS_TYPE_METERS SET_FOCUS_TYPE = common.FOCUS_TYPE_METERS
	// Focus automatically.
	FOCUS_TYPE_AUTO SET_FOCUS_TYPE = common.FOCUS_TYPE_AUTO
	// Single auto focus. Mainly used for still pictures. Usually abbreviated as AF-S.
	FOCUS_TYPE_AUTO_SINGLE SET_FOCUS_TYPE = common.FOCUS_TYPE_AUTO_SINGLE
	// Continuous auto focus. Mainly used for dynamic scenes. Abbreviated as AF-C.
	FOCUS_TYPE_AUTO_CONTINUOUS SET_FOCUS_TYPE = common.FOCUS_TYPE_AUTO_CONTINUOUS
)
