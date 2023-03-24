//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// A mapping of rover flight modes for custom_mode field of heartbeat.
type ROVER_MODE = ardupilotmega.ROVER_MODE

const (
	ROVER_MODE_MANUAL       ROVER_MODE = ardupilotmega.ROVER_MODE_MANUAL
	ROVER_MODE_ACRO         ROVER_MODE = ardupilotmega.ROVER_MODE_ACRO
	ROVER_MODE_STEERING     ROVER_MODE = ardupilotmega.ROVER_MODE_STEERING
	ROVER_MODE_HOLD         ROVER_MODE = ardupilotmega.ROVER_MODE_HOLD
	ROVER_MODE_LOITER       ROVER_MODE = ardupilotmega.ROVER_MODE_LOITER
	ROVER_MODE_FOLLOW       ROVER_MODE = ardupilotmega.ROVER_MODE_FOLLOW
	ROVER_MODE_SIMPLE       ROVER_MODE = ardupilotmega.ROVER_MODE_SIMPLE
	ROVER_MODE_AUTO         ROVER_MODE = ardupilotmega.ROVER_MODE_AUTO
	ROVER_MODE_RTL          ROVER_MODE = ardupilotmega.ROVER_MODE_RTL
	ROVER_MODE_SMART_RTL    ROVER_MODE = ardupilotmega.ROVER_MODE_SMART_RTL
	ROVER_MODE_GUIDED       ROVER_MODE = ardupilotmega.ROVER_MODE_GUIDED
	ROVER_MODE_INITIALIZING ROVER_MODE = ardupilotmega.ROVER_MODE_INITIALIZING
)
