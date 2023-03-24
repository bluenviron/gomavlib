//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package uavionix

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Enumeration of battery types
type MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE

const (
	// Not specified.
	MAV_BATTERY_TYPE_UNKNOWN MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE_UNKNOWN
	// Lithium polymer battery
	MAV_BATTERY_TYPE_LIPO MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE_LIPO
	// Lithium-iron-phosphate battery
	MAV_BATTERY_TYPE_LIFE MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE_LIFE
	// Lithium-ION battery
	MAV_BATTERY_TYPE_LION MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE_LION
	// Nickel metal hydride battery
	MAV_BATTERY_TYPE_NIMH MAV_BATTERY_TYPE = common.MAV_BATTERY_TYPE_NIMH
)
