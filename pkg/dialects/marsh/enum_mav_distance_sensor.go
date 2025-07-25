//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Enumeration of distance sensor types
type MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR

const (
	// Laser rangefinder, e.g. LightWare SF02/F or PulsedLight units
	MAV_DISTANCE_SENSOR_LASER MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR_LASER
	// Ultrasound rangefinder, e.g. MaxBotix units
	MAV_DISTANCE_SENSOR_ULTRASOUND MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR_ULTRASOUND
	// Infrared rangefinder, e.g. Sharp units
	MAV_DISTANCE_SENSOR_INFRARED MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR_INFRARED
	// Radar type, e.g. uLanding units
	MAV_DISTANCE_SENSOR_RADAR MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR_RADAR
	// Broken or unknown type, e.g. analog units
	MAV_DISTANCE_SENSOR_UNKNOWN MAV_DISTANCE_SENSOR = common.MAV_DISTANCE_SENSOR_UNKNOWN
)
