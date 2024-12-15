//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ualberta

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Camera tracking target data (shows where tracked target is within image)
type CAMERA_TRACKING_TARGET_DATA = common.CAMERA_TRACKING_TARGET_DATA

const (
	// No target data
	CAMERA_TRACKING_TARGET_DATA_NONE CAMERA_TRACKING_TARGET_DATA = common.CAMERA_TRACKING_TARGET_DATA_NONE
	// Target data embedded in image data (proprietary)
	CAMERA_TRACKING_TARGET_DATA_EMBEDDED CAMERA_TRACKING_TARGET_DATA = common.CAMERA_TRACKING_TARGET_DATA_EMBEDDED
	// Target data rendered in image
	CAMERA_TRACKING_TARGET_DATA_RENDERED CAMERA_TRACKING_TARGET_DATA = common.CAMERA_TRACKING_TARGET_DATA_RENDERED
	// Target data within status message (Point or Rectangle)
	CAMERA_TRACKING_TARGET_DATA_IN_STATUS CAMERA_TRACKING_TARGET_DATA = common.CAMERA_TRACKING_TARGET_DATA_IN_STATUS
)
