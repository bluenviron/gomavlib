//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package matrixpilot

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Camera tracking status flags
type CAMERA_TRACKING_STATUS_FLAGS = common.CAMERA_TRACKING_STATUS_FLAGS

const (
	// Camera is not tracking
	CAMERA_TRACKING_STATUS_FLAGS_IDLE CAMERA_TRACKING_STATUS_FLAGS = common.CAMERA_TRACKING_STATUS_FLAGS_IDLE
	// Camera is tracking
	CAMERA_TRACKING_STATUS_FLAGS_ACTIVE CAMERA_TRACKING_STATUS_FLAGS = common.CAMERA_TRACKING_STATUS_FLAGS_ACTIVE
	// Camera tracking in error state
	CAMERA_TRACKING_STATUS_FLAGS_ERROR CAMERA_TRACKING_STATUS_FLAGS = common.CAMERA_TRACKING_STATUS_FLAGS_ERROR
)
