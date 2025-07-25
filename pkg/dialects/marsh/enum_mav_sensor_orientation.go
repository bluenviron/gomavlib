//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Enumeration of sensor orientation, according to its rotations
type MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ORIENTATION

const (
	// Roll: 0, Pitch: 0, Yaw: 0
	MAV_SENSOR_ROTATION_NONE MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_NONE
	// Roll: 0, Pitch: 0, Yaw: 45
	MAV_SENSOR_ROTATION_YAW_45 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_45
	// Roll: 0, Pitch: 0, Yaw: 90
	MAV_SENSOR_ROTATION_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_90
	// Roll: 0, Pitch: 0, Yaw: 135
	MAV_SENSOR_ROTATION_YAW_135 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_135
	// Roll: 0, Pitch: 0, Yaw: 180
	MAV_SENSOR_ROTATION_YAW_180 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_180
	// Roll: 0, Pitch: 0, Yaw: 225
	MAV_SENSOR_ROTATION_YAW_225 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_225
	// Roll: 0, Pitch: 0, Yaw: 270
	MAV_SENSOR_ROTATION_YAW_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_270
	// Roll: 0, Pitch: 0, Yaw: 315
	MAV_SENSOR_ROTATION_YAW_315 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_YAW_315
	// Roll: 180, Pitch: 0, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_180 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180
	// Roll: 180, Pitch: 0, Yaw: 45
	MAV_SENSOR_ROTATION_ROLL_180_YAW_45 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_45
	// Roll: 180, Pitch: 0, Yaw: 90
	MAV_SENSOR_ROTATION_ROLL_180_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_90
	// Roll: 180, Pitch: 0, Yaw: 135
	MAV_SENSOR_ROTATION_ROLL_180_YAW_135 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_135
	// Roll: 0, Pitch: 180, Yaw: 0
	MAV_SENSOR_ROTATION_PITCH_180 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_180
	// Roll: 180, Pitch: 0, Yaw: 225
	MAV_SENSOR_ROTATION_ROLL_180_YAW_225 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_225
	// Roll: 180, Pitch: 0, Yaw: 270
	MAV_SENSOR_ROTATION_ROLL_180_YAW_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_270
	// Roll: 180, Pitch: 0, Yaw: 315
	MAV_SENSOR_ROTATION_ROLL_180_YAW_315 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_YAW_315
	// Roll: 90, Pitch: 0, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90
	// Roll: 90, Pitch: 0, Yaw: 45
	MAV_SENSOR_ROTATION_ROLL_90_YAW_45 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_YAW_45
	// Roll: 90, Pitch: 0, Yaw: 90
	MAV_SENSOR_ROTATION_ROLL_90_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_YAW_90
	// Roll: 90, Pitch: 0, Yaw: 135
	MAV_SENSOR_ROTATION_ROLL_90_YAW_135 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_YAW_135
	// Roll: 270, Pitch: 0, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270
	// Roll: 270, Pitch: 0, Yaw: 45
	MAV_SENSOR_ROTATION_ROLL_270_YAW_45 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_YAW_45
	// Roll: 270, Pitch: 0, Yaw: 90
	MAV_SENSOR_ROTATION_ROLL_270_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_YAW_90
	// Roll: 270, Pitch: 0, Yaw: 135
	MAV_SENSOR_ROTATION_ROLL_270_YAW_135 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_YAW_135
	// Roll: 0, Pitch: 90, Yaw: 0
	MAV_SENSOR_ROTATION_PITCH_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_90
	// Roll: 0, Pitch: 270, Yaw: 0
	MAV_SENSOR_ROTATION_PITCH_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_270
	// Roll: 0, Pitch: 180, Yaw: 90
	MAV_SENSOR_ROTATION_PITCH_180_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_180_YAW_90
	// Roll: 0, Pitch: 180, Yaw: 270
	MAV_SENSOR_ROTATION_PITCH_180_YAW_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_180_YAW_270
	// Roll: 90, Pitch: 90, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_90
	// Roll: 180, Pitch: 90, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_180_PITCH_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_PITCH_90
	// Roll: 270, Pitch: 90, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_270_PITCH_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_PITCH_90
	// Roll: 90, Pitch: 180, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_180 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_180
	// Roll: 270, Pitch: 180, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_270_PITCH_180 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_PITCH_180
	// Roll: 90, Pitch: 270, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_270
	// Roll: 180, Pitch: 270, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_180_PITCH_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_180_PITCH_270
	// Roll: 270, Pitch: 270, Yaw: 0
	MAV_SENSOR_ROTATION_ROLL_270_PITCH_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_270_PITCH_270
	// Roll: 90, Pitch: 180, Yaw: 90
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_180_YAW_90 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_180_YAW_90
	// Roll: 90, Pitch: 0, Yaw: 270
	MAV_SENSOR_ROTATION_ROLL_90_YAW_270 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_YAW_270
	// Roll: 90, Pitch: 68, Yaw: 293
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_68_YAW_293 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_68_YAW_293
	// Pitch: 315
	MAV_SENSOR_ROTATION_PITCH_315 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_PITCH_315
	// Roll: 90, Pitch: 315
	MAV_SENSOR_ROTATION_ROLL_90_PITCH_315 MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_ROLL_90_PITCH_315
	// Custom orientation
	MAV_SENSOR_ROTATION_CUSTOM MAV_SENSOR_ORIENTATION = common.MAV_SENSOR_ROTATION_CUSTOM
)
