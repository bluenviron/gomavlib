//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Actuator configuration, used to change a setting on an actuator. Component information metadata can be used to know which outputs support which commands.
type ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION

const (
	// Do nothing.
	ACTUATOR_CONFIGURATION_NONE ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_NONE
	// Command the actuator to beep now.
	ACTUATOR_CONFIGURATION_BEEP ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_BEEP
	// Permanently set the actuator (ESC) to 3D mode (reversible thrust).
	ACTUATOR_CONFIGURATION_3D_MODE_ON ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_3D_MODE_ON
	// Permanently set the actuator (ESC) to non 3D mode (non-reversible thrust).
	ACTUATOR_CONFIGURATION_3D_MODE_OFF ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_3D_MODE_OFF
	// Permanently set the actuator (ESC) to spin direction 1 (which can be clockwise or counter-clockwise).
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION1 ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_SPIN_DIRECTION1
	// Permanently set the actuator (ESC) to spin direction 2 (opposite of direction 1).
	ACTUATOR_CONFIGURATION_SPIN_DIRECTION2 ACTUATOR_CONFIGURATION = common.ACTUATOR_CONFIGURATION_SPIN_DIRECTION2
)
