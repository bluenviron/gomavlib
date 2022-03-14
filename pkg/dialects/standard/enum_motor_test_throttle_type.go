//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl
package standard

import (
	"errors"
)

// Defines how throttle value is represented in MAV_CMD_DO_MOTOR_TEST.
type MOTOR_TEST_THROTTLE_TYPE uint32

const (
	// Throttle as a percentage (0 ~ 100)
	MOTOR_TEST_THROTTLE_PERCENT MOTOR_TEST_THROTTLE_TYPE = 0
	// Throttle as an absolute PWM value (normally in range of 1000~2000).
	MOTOR_TEST_THROTTLE_PWM MOTOR_TEST_THROTTLE_TYPE = 1
	// Throttle pass-through from pilot's transmitter.
	MOTOR_TEST_THROTTLE_PILOT MOTOR_TEST_THROTTLE_TYPE = 2
	// Per-motor compass calibration test.
	MOTOR_TEST_COMPASS_CAL MOTOR_TEST_THROTTLE_TYPE = 3
)

var labels_MOTOR_TEST_THROTTLE_TYPE = map[MOTOR_TEST_THROTTLE_TYPE]string{}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MOTOR_TEST_THROTTLE_TYPE) MarshalText() ([]byte, error) {
	if l, ok := labels_MOTOR_TEST_THROTTLE_TYPE[e]; ok {
		return []byte(l), nil
	}
	return nil, errors.New("invalid value")
}

var reverseLabels_MOTOR_TEST_THROTTLE_TYPE = map[string]MOTOR_TEST_THROTTLE_TYPE{}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MOTOR_TEST_THROTTLE_TYPE) UnmarshalText(text []byte) error {
	if rl, ok := reverseLabels_MOTOR_TEST_THROTTLE_TYPE[string(text)]; ok {
		*e = rl
		return nil
	}
	return errors.New("invalid value")
}

// String implements the fmt.Stringer interface.
func (e MOTOR_TEST_THROTTLE_TYPE) String() string {
	if l, ok := labels_MOTOR_TEST_THROTTLE_TYPE[e]; ok {
		return l
	}
	return "invalid value"
}