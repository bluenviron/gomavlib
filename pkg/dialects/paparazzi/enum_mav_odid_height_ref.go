//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

type MAV_ODID_HEIGHT_REF = common.MAV_ODID_HEIGHT_REF

const (
	// The height field is relative to the take-off location.
	MAV_ODID_HEIGHT_REF_OVER_TAKEOFF MAV_ODID_HEIGHT_REF = common.MAV_ODID_HEIGHT_REF_OVER_TAKEOFF
	// The height field is relative to ground.
	MAV_ODID_HEIGHT_REF_OVER_GROUND MAV_ODID_HEIGHT_REF = common.MAV_ODID_HEIGHT_REF_OVER_GROUND
)
