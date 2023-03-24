//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ualberta

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Parachute actions. Trigger release and enable/disable auto-release.
type PARACHUTE_ACTION = common.PARACHUTE_ACTION

const (
	// Disable auto-release of parachute (i.e. release triggered by crash detectors).
	PARACHUTE_DISABLE PARACHUTE_ACTION = common.PARACHUTE_DISABLE
	// Enable auto-release of parachute.
	PARACHUTE_ENABLE PARACHUTE_ACTION = common.PARACHUTE_ENABLE
	// Release parachute and kill motors.
	PARACHUTE_RELEASE PARACHUTE_ACTION = common.PARACHUTE_RELEASE
)
