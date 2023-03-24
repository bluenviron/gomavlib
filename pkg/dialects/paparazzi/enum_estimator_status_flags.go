//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Flags in ESTIMATOR_STATUS message
type ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_STATUS_FLAGS

const (
	// True if the attitude estimate is good
	ESTIMATOR_ATTITUDE ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_ATTITUDE
	// True if the horizontal velocity estimate is good
	ESTIMATOR_VELOCITY_HORIZ ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_VELOCITY_HORIZ
	// True if the  vertical velocity estimate is good
	ESTIMATOR_VELOCITY_VERT ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_VELOCITY_VERT
	// True if the horizontal position (relative) estimate is good
	ESTIMATOR_POS_HORIZ_REL ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_POS_HORIZ_REL
	// True if the horizontal position (absolute) estimate is good
	ESTIMATOR_POS_HORIZ_ABS ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_POS_HORIZ_ABS
	// True if the vertical position (absolute) estimate is good
	ESTIMATOR_POS_VERT_ABS ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_POS_VERT_ABS
	// True if the vertical position (above ground) estimate is good
	ESTIMATOR_POS_VERT_AGL ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_POS_VERT_AGL
	// True if the EKF is in a constant position mode and is not using external measurements (eg GPS or optical flow)
	ESTIMATOR_CONST_POS_MODE ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_CONST_POS_MODE
	// True if the EKF has sufficient data to enter a mode that will provide a (relative) position estimate
	ESTIMATOR_PRED_POS_HORIZ_REL ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_PRED_POS_HORIZ_REL
	// True if the EKF has sufficient data to enter a mode that will provide a (absolute) position estimate
	ESTIMATOR_PRED_POS_HORIZ_ABS ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_PRED_POS_HORIZ_ABS
	// True if the EKF has detected a GPS glitch
	ESTIMATOR_GPS_GLITCH ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_GPS_GLITCH
	// True if the EKF has detected bad accelerometer data
	ESTIMATOR_ACCEL_ERROR ESTIMATOR_STATUS_FLAGS = common.ESTIMATOR_ACCEL_ERROR
)
