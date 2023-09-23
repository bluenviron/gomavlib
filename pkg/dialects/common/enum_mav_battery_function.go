//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Enumeration of battery functions
type MAV_BATTERY_FUNCTION uint32

const (
	// Battery function is unknown
	MAV_BATTERY_FUNCTION_UNKNOWN MAV_BATTERY_FUNCTION = 0
	// Battery supports all flight systems
	MAV_BATTERY_FUNCTION_ALL MAV_BATTERY_FUNCTION = 1
	// Battery for the propulsion system
	MAV_BATTERY_FUNCTION_PROPULSION MAV_BATTERY_FUNCTION = 2
	// Avionics battery
	MAV_BATTERY_FUNCTION_AVIONICS MAV_BATTERY_FUNCTION = 3
	// Payload battery
	MAV_BATTERY_FUNCTION_PAYLOAD MAV_BATTERY_FUNCTION = 4
)

var labels_MAV_BATTERY_FUNCTION = map[MAV_BATTERY_FUNCTION]string{
	MAV_BATTERY_FUNCTION_UNKNOWN:    "MAV_BATTERY_FUNCTION_UNKNOWN",
	MAV_BATTERY_FUNCTION_ALL:        "MAV_BATTERY_FUNCTION_ALL",
	MAV_BATTERY_FUNCTION_PROPULSION: "MAV_BATTERY_FUNCTION_PROPULSION",
	MAV_BATTERY_FUNCTION_AVIONICS:   "MAV_BATTERY_FUNCTION_AVIONICS",
	MAV_BATTERY_FUNCTION_PAYLOAD:    "MAV_BATTERY_FUNCTION_PAYLOAD",
}

var values_MAV_BATTERY_FUNCTION = map[string]MAV_BATTERY_FUNCTION{
	"MAV_BATTERY_FUNCTION_UNKNOWN":    MAV_BATTERY_FUNCTION_UNKNOWN,
	"MAV_BATTERY_FUNCTION_ALL":        MAV_BATTERY_FUNCTION_ALL,
	"MAV_BATTERY_FUNCTION_PROPULSION": MAV_BATTERY_FUNCTION_PROPULSION,
	"MAV_BATTERY_FUNCTION_AVIONICS":   MAV_BATTERY_FUNCTION_AVIONICS,
	"MAV_BATTERY_FUNCTION_PAYLOAD":    MAV_BATTERY_FUNCTION_PAYLOAD,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_BATTERY_FUNCTION) MarshalText() ([]byte, error) {
	name, ok := labels_MAV_BATTERY_FUNCTION[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_BATTERY_FUNCTION) UnmarshalText(text []byte) error {
	value, ok := values_MAV_BATTERY_FUNCTION[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_BATTERY_FUNCTION) String() string {
	name, ok := labels_MAV_BATTERY_FUNCTION[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
