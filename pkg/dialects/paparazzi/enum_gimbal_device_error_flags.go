//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Gimbal device (low level) error flags (bitmap, 0 means no error)
type GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS

const (
	// Gimbal device is limited by hardware roll limit.
	GIMBAL_DEVICE_ERROR_FLAGS_AT_ROLL_LIMIT GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_AT_ROLL_LIMIT
	// Gimbal device is limited by hardware pitch limit.
	GIMBAL_DEVICE_ERROR_FLAGS_AT_PITCH_LIMIT GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_AT_PITCH_LIMIT
	// Gimbal device is limited by hardware yaw limit.
	GIMBAL_DEVICE_ERROR_FLAGS_AT_YAW_LIMIT GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_AT_YAW_LIMIT
	// There is an error with the gimbal encoders.
	GIMBAL_DEVICE_ERROR_FLAGS_ENCODER_ERROR GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_ENCODER_ERROR
	// There is an error with the gimbal power source.
	GIMBAL_DEVICE_ERROR_FLAGS_POWER_ERROR GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_POWER_ERROR
	// There is an error with the gimbal motors.
	GIMBAL_DEVICE_ERROR_FLAGS_MOTOR_ERROR GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_MOTOR_ERROR
	// There is an error with the gimbal's software.
	GIMBAL_DEVICE_ERROR_FLAGS_SOFTWARE_ERROR GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_SOFTWARE_ERROR
	// There is an error with the gimbal's communication.
	GIMBAL_DEVICE_ERROR_FLAGS_COMMS_ERROR GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_COMMS_ERROR
	// Gimbal device is currently calibrating.
	GIMBAL_DEVICE_ERROR_FLAGS_CALIBRATION_RUNNING GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_CALIBRATION_RUNNING
	// Gimbal device is not assigned to a gimbal manager.
	GIMBAL_DEVICE_ERROR_FLAGS_NO_MANAGER GIMBAL_DEVICE_ERROR_FLAGS = common.GIMBAL_DEVICE_ERROR_FLAGS_NO_MANAGER
)
