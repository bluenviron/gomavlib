//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Enumeration of landed detector states
type MAV_LANDED_STATE uint32

const (
	// MAV landed state is unknown
	MAV_LANDED_STATE_UNDEFINED MAV_LANDED_STATE = 0
	// MAV is landed (on ground)
	MAV_LANDED_STATE_ON_GROUND MAV_LANDED_STATE = 1
	// MAV is in air
	MAV_LANDED_STATE_IN_AIR MAV_LANDED_STATE = 2
	// MAV currently taking off
	MAV_LANDED_STATE_TAKEOFF MAV_LANDED_STATE = 3
	// MAV currently landing
	MAV_LANDED_STATE_LANDING MAV_LANDED_STATE = 4
)

var labels_MAV_LANDED_STATE = map[MAV_LANDED_STATE]string{
	MAV_LANDED_STATE_UNDEFINED: "MAV_LANDED_STATE_UNDEFINED",
	MAV_LANDED_STATE_ON_GROUND: "MAV_LANDED_STATE_ON_GROUND",
	MAV_LANDED_STATE_IN_AIR:    "MAV_LANDED_STATE_IN_AIR",
	MAV_LANDED_STATE_TAKEOFF:   "MAV_LANDED_STATE_TAKEOFF",
	MAV_LANDED_STATE_LANDING:   "MAV_LANDED_STATE_LANDING",
}

var values_MAV_LANDED_STATE = map[string]MAV_LANDED_STATE{
	"MAV_LANDED_STATE_UNDEFINED": MAV_LANDED_STATE_UNDEFINED,
	"MAV_LANDED_STATE_ON_GROUND": MAV_LANDED_STATE_ON_GROUND,
	"MAV_LANDED_STATE_IN_AIR":    MAV_LANDED_STATE_IN_AIR,
	"MAV_LANDED_STATE_TAKEOFF":   MAV_LANDED_STATE_TAKEOFF,
	"MAV_LANDED_STATE_LANDING":   MAV_LANDED_STATE_LANDING,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_LANDED_STATE) MarshalText() ([]byte, error) {
	if name, ok := labels_MAV_LANDED_STATE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_LANDED_STATE) UnmarshalText(text []byte) error {
	if value, ok := values_MAV_LANDED_STATE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = MAV_LANDED_STATE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_LANDED_STATE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
