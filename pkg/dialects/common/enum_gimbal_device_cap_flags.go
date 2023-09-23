//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strings"
)

// Gimbal device (low level) capability flags (bitmap).
type GIMBAL_DEVICE_CAP_FLAGS uint32

const (
	// Gimbal device supports a retracted position.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_RETRACT GIMBAL_DEVICE_CAP_FLAGS = 1
	// Gimbal device supports a horizontal, forward looking position, stabilized.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_NEUTRAL GIMBAL_DEVICE_CAP_FLAGS = 2
	// Gimbal device supports rotating around roll axis.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_AXIS GIMBAL_DEVICE_CAP_FLAGS = 4
	// Gimbal device supports to follow a roll angle relative to the vehicle.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_FOLLOW GIMBAL_DEVICE_CAP_FLAGS = 8
	// Gimbal device supports locking to a roll angle (generally that's the default with roll stabilized).
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_LOCK GIMBAL_DEVICE_CAP_FLAGS = 16
	// Gimbal device supports rotating around pitch axis.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_AXIS GIMBAL_DEVICE_CAP_FLAGS = 32
	// Gimbal device supports to follow a pitch angle relative to the vehicle.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_FOLLOW GIMBAL_DEVICE_CAP_FLAGS = 64
	// Gimbal device supports locking to a pitch angle (generally that's the default with pitch stabilized).
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_LOCK GIMBAL_DEVICE_CAP_FLAGS = 128
	// Gimbal device supports rotating around yaw axis.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_AXIS GIMBAL_DEVICE_CAP_FLAGS = 256
	// Gimbal device supports to follow a yaw angle relative to the vehicle (generally that's the default).
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_FOLLOW GIMBAL_DEVICE_CAP_FLAGS = 512
	// Gimbal device supports locking to an absolute heading, i.e., yaw angle relative to North (earth frame, often this is an option available).
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_LOCK GIMBAL_DEVICE_CAP_FLAGS = 1024
	// Gimbal device supports yawing/panning infinetely (e.g. using slip disk).
	GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_INFINITE_YAW GIMBAL_DEVICE_CAP_FLAGS = 2048
	// Gimbal device supports yaw angles and angular velocities relative to North (earth frame). This usually requires support by an autopilot via AUTOPILOT_STATE_FOR_GIMBAL_DEVICE. Support can go on and off during runtime, which is reported by the flag GIMBAL_DEVICE_FLAGS_CAN_ACCEPT_YAW_IN_EARTH_FRAME.
	GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_YAW_IN_EARTH_FRAME GIMBAL_DEVICE_CAP_FLAGS = 4096
	// Gimbal device supports radio control inputs as an alternative input for controlling the gimbal orientation.
	GIMBAL_DEVICE_CAP_FLAGS_HAS_RC_INPUTS GIMBAL_DEVICE_CAP_FLAGS = 8192
)

var labels_GIMBAL_DEVICE_CAP_FLAGS = map[GIMBAL_DEVICE_CAP_FLAGS]string{
	GIMBAL_DEVICE_CAP_FLAGS_HAS_RETRACT:                 "GIMBAL_DEVICE_CAP_FLAGS_HAS_RETRACT",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_NEUTRAL:                 "GIMBAL_DEVICE_CAP_FLAGS_HAS_NEUTRAL",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_AXIS:               "GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_AXIS",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_FOLLOW:             "GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_FOLLOW",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_LOCK:               "GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_LOCK",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_AXIS:              "GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_AXIS",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_FOLLOW:            "GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_FOLLOW",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_LOCK:              "GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_LOCK",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_AXIS:                "GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_AXIS",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_FOLLOW:              "GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_FOLLOW",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_LOCK:                "GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_LOCK",
	GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_INFINITE_YAW:       "GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_INFINITE_YAW",
	GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_YAW_IN_EARTH_FRAME: "GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_YAW_IN_EARTH_FRAME",
	GIMBAL_DEVICE_CAP_FLAGS_HAS_RC_INPUTS:               "GIMBAL_DEVICE_CAP_FLAGS_HAS_RC_INPUTS",
}

var values_GIMBAL_DEVICE_CAP_FLAGS = map[string]GIMBAL_DEVICE_CAP_FLAGS{
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_RETRACT":                 GIMBAL_DEVICE_CAP_FLAGS_HAS_RETRACT,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_NEUTRAL":                 GIMBAL_DEVICE_CAP_FLAGS_HAS_NEUTRAL,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_AXIS":               GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_AXIS,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_FOLLOW":             GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_FOLLOW,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_LOCK":               GIMBAL_DEVICE_CAP_FLAGS_HAS_ROLL_LOCK,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_AXIS":              GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_AXIS,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_FOLLOW":            GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_FOLLOW,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_LOCK":              GIMBAL_DEVICE_CAP_FLAGS_HAS_PITCH_LOCK,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_AXIS":                GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_AXIS,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_FOLLOW":              GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_FOLLOW,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_LOCK":                GIMBAL_DEVICE_CAP_FLAGS_HAS_YAW_LOCK,
	"GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_INFINITE_YAW":       GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_INFINITE_YAW,
	"GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_YAW_IN_EARTH_FRAME": GIMBAL_DEVICE_CAP_FLAGS_SUPPORTS_YAW_IN_EARTH_FRAME,
	"GIMBAL_DEVICE_CAP_FLAGS_HAS_RC_INPUTS":               GIMBAL_DEVICE_CAP_FLAGS_HAS_RC_INPUTS,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GIMBAL_DEVICE_CAP_FLAGS) MarshalText() ([]byte, error) {
	var names []string
	for i := 0; i < 14; i++ {
		mask := GIMBAL_DEVICE_CAP_FLAGS(1 << i)
		if e&mask == mask {
			names = append(names, labels_GIMBAL_DEVICE_CAP_FLAGS[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GIMBAL_DEVICE_CAP_FLAGS) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask GIMBAL_DEVICE_CAP_FLAGS
	for _, label := range labels {
		if value, ok := values_GIMBAL_DEVICE_CAP_FLAGS[label]; ok {
			mask |= value
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e GIMBAL_DEVICE_CAP_FLAGS) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
