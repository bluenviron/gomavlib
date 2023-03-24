//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Battery mode. Note, the normal operation mode (i.e. when flying) should be reported as MAV_BATTERY_MODE_UNKNOWN to allow message trimming in normal flight.
type MAV_BATTERY_MODE = common.MAV_BATTERY_MODE

const (
	// Battery mode not supported/unknown battery mode/normal operation.
	MAV_BATTERY_MODE_UNKNOWN MAV_BATTERY_MODE = common.MAV_BATTERY_MODE_UNKNOWN
	// Battery is auto discharging (towards storage level).
	MAV_BATTERY_MODE_AUTO_DISCHARGING MAV_BATTERY_MODE = common.MAV_BATTERY_MODE_AUTO_DISCHARGING
	// Battery in hot-swap mode (current limited to prevent spikes that might damage sensitive electrical circuits).
	MAV_BATTERY_MODE_HOT_SWAP MAV_BATTERY_MODE = common.MAV_BATTERY_MODE_HOT_SWAP
)
