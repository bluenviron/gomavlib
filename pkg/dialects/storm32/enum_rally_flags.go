//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// Flags in RALLY_POINT message.
type RALLY_FLAGS = ardupilotmega.RALLY_FLAGS

const (
	// Flag set when requiring favorable winds for landing.
	FAVORABLE_WIND RALLY_FLAGS = ardupilotmega.FAVORABLE_WIND
	// Flag set when plane is to immediately descend to break altitude and land without GCS intervention. Flag not set when plane is to loiter at Rally point until commanded to land.
	LAND_IMMEDIATELY RALLY_FLAGS = ardupilotmega.LAND_IMMEDIATELY
)
