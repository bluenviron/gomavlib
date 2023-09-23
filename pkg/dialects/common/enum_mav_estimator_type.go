//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Enumeration of estimator types
type MAV_ESTIMATOR_TYPE uint32

const (
	// Unknown type of the estimator.
	MAV_ESTIMATOR_TYPE_UNKNOWN MAV_ESTIMATOR_TYPE = 0
	// This is a naive estimator without any real covariance feedback.
	MAV_ESTIMATOR_TYPE_NAIVE MAV_ESTIMATOR_TYPE = 1
	// Computer vision based estimate. Might be up to scale.
	MAV_ESTIMATOR_TYPE_VISION MAV_ESTIMATOR_TYPE = 2
	// Visual-inertial estimate.
	MAV_ESTIMATOR_TYPE_VIO MAV_ESTIMATOR_TYPE = 3
	// Plain GPS estimate.
	MAV_ESTIMATOR_TYPE_GPS MAV_ESTIMATOR_TYPE = 4
	// Estimator integrating GPS and inertial sensing.
	MAV_ESTIMATOR_TYPE_GPS_INS MAV_ESTIMATOR_TYPE = 5
	// Estimate from external motion capturing system.
	MAV_ESTIMATOR_TYPE_MOCAP MAV_ESTIMATOR_TYPE = 6
	// Estimator based on lidar sensor input.
	MAV_ESTIMATOR_TYPE_LIDAR MAV_ESTIMATOR_TYPE = 7
	// Estimator on autopilot.
	MAV_ESTIMATOR_TYPE_AUTOPILOT MAV_ESTIMATOR_TYPE = 8
)

var labels_MAV_ESTIMATOR_TYPE = map[MAV_ESTIMATOR_TYPE]string{
	MAV_ESTIMATOR_TYPE_UNKNOWN:   "MAV_ESTIMATOR_TYPE_UNKNOWN",
	MAV_ESTIMATOR_TYPE_NAIVE:     "MAV_ESTIMATOR_TYPE_NAIVE",
	MAV_ESTIMATOR_TYPE_VISION:    "MAV_ESTIMATOR_TYPE_VISION",
	MAV_ESTIMATOR_TYPE_VIO:       "MAV_ESTIMATOR_TYPE_VIO",
	MAV_ESTIMATOR_TYPE_GPS:       "MAV_ESTIMATOR_TYPE_GPS",
	MAV_ESTIMATOR_TYPE_GPS_INS:   "MAV_ESTIMATOR_TYPE_GPS_INS",
	MAV_ESTIMATOR_TYPE_MOCAP:     "MAV_ESTIMATOR_TYPE_MOCAP",
	MAV_ESTIMATOR_TYPE_LIDAR:     "MAV_ESTIMATOR_TYPE_LIDAR",
	MAV_ESTIMATOR_TYPE_AUTOPILOT: "MAV_ESTIMATOR_TYPE_AUTOPILOT",
}

var values_MAV_ESTIMATOR_TYPE = map[string]MAV_ESTIMATOR_TYPE{
	"MAV_ESTIMATOR_TYPE_UNKNOWN":   MAV_ESTIMATOR_TYPE_UNKNOWN,
	"MAV_ESTIMATOR_TYPE_NAIVE":     MAV_ESTIMATOR_TYPE_NAIVE,
	"MAV_ESTIMATOR_TYPE_VISION":    MAV_ESTIMATOR_TYPE_VISION,
	"MAV_ESTIMATOR_TYPE_VIO":       MAV_ESTIMATOR_TYPE_VIO,
	"MAV_ESTIMATOR_TYPE_GPS":       MAV_ESTIMATOR_TYPE_GPS,
	"MAV_ESTIMATOR_TYPE_GPS_INS":   MAV_ESTIMATOR_TYPE_GPS_INS,
	"MAV_ESTIMATOR_TYPE_MOCAP":     MAV_ESTIMATOR_TYPE_MOCAP,
	"MAV_ESTIMATOR_TYPE_LIDAR":     MAV_ESTIMATOR_TYPE_LIDAR,
	"MAV_ESTIMATOR_TYPE_AUTOPILOT": MAV_ESTIMATOR_TYPE_AUTOPILOT,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_ESTIMATOR_TYPE) MarshalText() ([]byte, error) {
	name, ok := labels_MAV_ESTIMATOR_TYPE[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_ESTIMATOR_TYPE) UnmarshalText(text []byte) error {
	value, ok := values_MAV_ESTIMATOR_TYPE[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_ESTIMATOR_TYPE) String() string {
	name, ok := labels_MAV_ESTIMATOR_TYPE[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
