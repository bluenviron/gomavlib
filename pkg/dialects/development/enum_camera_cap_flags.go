//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package development

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Camera capability flags (Bitmap)
type CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS

const (
	// Camera is able to record video
	CAMERA_CAP_FLAGS_CAPTURE_VIDEO CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_CAPTURE_VIDEO
	// Camera is able to capture images
	CAMERA_CAP_FLAGS_CAPTURE_IMAGE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_CAPTURE_IMAGE
	// Camera has separate Video and Image/Photo modes (MAV_CMD_SET_CAMERA_MODE)
	CAMERA_CAP_FLAGS_HAS_MODES CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_MODES
	// Camera can capture images while in video mode
	CAMERA_CAP_FLAGS_CAN_CAPTURE_IMAGE_IN_VIDEO_MODE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_CAN_CAPTURE_IMAGE_IN_VIDEO_MODE
	// Camera can capture videos while in Photo/Image mode
	CAMERA_CAP_FLAGS_CAN_CAPTURE_VIDEO_IN_IMAGE_MODE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_CAN_CAPTURE_VIDEO_IN_IMAGE_MODE
	// Camera has image survey mode (MAV_CMD_SET_CAMERA_MODE)
	CAMERA_CAP_FLAGS_HAS_IMAGE_SURVEY_MODE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_IMAGE_SURVEY_MODE
	// Camera has basic zoom control (MAV_CMD_SET_CAMERA_ZOOM)
	CAMERA_CAP_FLAGS_HAS_BASIC_ZOOM CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_BASIC_ZOOM
	// Camera has basic focus control (MAV_CMD_SET_CAMERA_FOCUS)
	CAMERA_CAP_FLAGS_HAS_BASIC_FOCUS CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_BASIC_FOCUS
	// Camera has video streaming capabilities (request VIDEO_STREAM_INFORMATION with MAV_CMD_REQUEST_MESSAGE for video streaming info)
	CAMERA_CAP_FLAGS_HAS_VIDEO_STREAM CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_VIDEO_STREAM
	// Camera supports tracking of a point on the camera view.
	CAMERA_CAP_FLAGS_HAS_TRACKING_POINT CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_TRACKING_POINT
	// Camera supports tracking of a selection rectangle on the camera view.
	CAMERA_CAP_FLAGS_HAS_TRACKING_RECTANGLE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_TRACKING_RECTANGLE
	// Camera supports tracking geo status (CAMERA_TRACKING_GEO_STATUS).
	CAMERA_CAP_FLAGS_HAS_TRACKING_GEO_STATUS CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_TRACKING_GEO_STATUS
	// Camera supports absolute thermal range (request CAMERA_THERMAL_RANGE with MAV_CMD_REQUEST_MESSAGE) (WIP).
	CAMERA_CAP_FLAGS_HAS_THERMAL_RANGE CAMERA_CAP_FLAGS = common.CAMERA_CAP_FLAGS_HAS_THERMAL_RANGE
)
