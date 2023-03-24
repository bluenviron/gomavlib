//autogenerated:yes
//nolint:revive,misspell,govet,lll
package pythonarraytest

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Low level message to control a gimbal device's attitude.
// This message is to be sent from the gimbal manager to the gimbal device component.
// The quaternion and angular velocities can be set to NaN according to use case.
// For the angles encoded in the quaternion and the angular velocities holds:
// If the flag GIMBAL_DEVICE_FLAGS_YAW_IN_VEHICLE_FRAME is set, then they are relative to the vehicle heading (vehicle frame).
// If the flag GIMBAL_DEVICE_FLAGS_YAW_IN_EARTH_FRAME is set, then they are relative to absolute North (earth frame).
// If neither of these flags are set, then (for backwards compatibility) it holds:
// If the flag GIMBAL_DEVICE_FLAGS_YAW_LOCK is set, then they are relative to absolute North (earth frame),
// else they are relative to the vehicle heading (vehicle frame).
// Setting both GIMBAL_DEVICE_FLAGS_YAW_IN_VEHICLE_FRAME and GIMBAL_DEVICE_FLAGS_YAW_IN_EARTH_FRAME is not allowed.
// These rules are to ensure backwards compatibility.
// New implementations should always set either GIMBAL_DEVICE_FLAGS_YAW_IN_VEHICLE_FRAME or GIMBAL_DEVICE_FLAGS_YAW_IN_EARTH_FRAME.
type MessageGimbalDeviceSetAttitude = common.MessageGimbalDeviceSetAttitude
