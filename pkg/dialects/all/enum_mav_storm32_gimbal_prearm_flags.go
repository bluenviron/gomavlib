//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/storm32"
)

// STorM32 gimbal prearm check flags.
type MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS

const (
	// STorM32 gimbal is in normal state.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_IS_NORMAL MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_IS_NORMAL
	// The IMUs are healthy and working normally.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_IMUS_WORKING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_IMUS_WORKING
	// The motors are active and working normally.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_MOTORS_WORKING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_MOTORS_WORKING
	// The encoders are healthy and working normally.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_ENCODERS_WORKING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_ENCODERS_WORKING
	// A battery voltage is applied and is in range.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_VOLTAGE_OK MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_VOLTAGE_OK
	// Virtual input channels are receiving data.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_VIRTUALCHANNELS_RECEIVING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_VIRTUALCHANNELS_RECEIVING
	// Mavlink messages are being received.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_MAVLINK_RECEIVING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_MAVLINK_RECEIVING
	// The STorM32Link data indicates QFix.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_STORM32LINK_QFIX MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_STORM32LINK_QFIX
	// The STorM32Link is working.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_STORM32LINK_WORKING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_STORM32LINK_WORKING
	// The camera has been found and is connected.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_CAMERA_CONNECTED MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_CAMERA_CONNECTED
	// The signal on the AUX0 input pin is low.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_AUX0_LOW MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_AUX0_LOW
	// The signal on the AUX1 input pin is low.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_AUX1_LOW MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_AUX1_LOW
	// The NTLogger is working normally.
	MAV_STORM32_GIMBAL_PREARM_FLAGS_NTLOGGER_WORKING MAV_STORM32_GIMBAL_PREARM_FLAGS = storm32.MAV_STORM32_GIMBAL_PREARM_FLAGS_NTLOGGER_WORKING
)
