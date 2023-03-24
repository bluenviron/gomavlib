//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Sequence that motors are tested when using MAV_CMD_DO_MOTOR_TEST.
type MOTOR_TEST_ORDER = common.MOTOR_TEST_ORDER

const (
	// Default autopilot motor test method.
	MOTOR_TEST_ORDER_DEFAULT MOTOR_TEST_ORDER = common.MOTOR_TEST_ORDER_DEFAULT
	// Motor numbers are specified as their index in a predefined vehicle-specific sequence.
	MOTOR_TEST_ORDER_SEQUENCE MOTOR_TEST_ORDER = common.MOTOR_TEST_ORDER_SEQUENCE
	// Motor numbers are specified as the output as labeled on the board.
	MOTOR_TEST_ORDER_BOARD MOTOR_TEST_ORDER = common.MOTOR_TEST_ORDER_BOARD
)
