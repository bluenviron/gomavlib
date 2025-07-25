//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Camera tracking modes
type CAMERA_TRACKING_MODE = common.CAMERA_TRACKING_MODE

const (
	// Not tracking
	CAMERA_TRACKING_MODE_NONE CAMERA_TRACKING_MODE = common.CAMERA_TRACKING_MODE_NONE
	// Target is a point
	CAMERA_TRACKING_MODE_POINT CAMERA_TRACKING_MODE = common.CAMERA_TRACKING_MODE_POINT
	// Target is a rectangle
	CAMERA_TRACKING_MODE_RECTANGLE CAMERA_TRACKING_MODE = common.CAMERA_TRACKING_MODE_RECTANGLE
)
