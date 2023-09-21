//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package minimal

import (
	"fmt"
	"strconv"
)

type MAV_STATE uint32

const (
	// Uninitialized system, state is unknown.
	MAV_STATE_UNINIT MAV_STATE = 0
	// System is booting up.
	MAV_STATE_BOOT MAV_STATE = 1
	// System is calibrating and not flight-ready.
	MAV_STATE_CALIBRATING MAV_STATE = 2
	// System is grounded and on standby. It can be launched any time.
	MAV_STATE_STANDBY MAV_STATE = 3
	// System is active and might be already airborne. Motors are engaged.
	MAV_STATE_ACTIVE MAV_STATE = 4
	// System is in a non-normal flight mode (failsafe). It can however still navigate.
	MAV_STATE_CRITICAL MAV_STATE = 5
	// System is in a non-normal flight mode (failsafe). It lost control over parts or over the whole airframe. It is in mayday and going down.
	MAV_STATE_EMERGENCY MAV_STATE = 6
	// System just initialized its power-down sequence, will shut down now.
	MAV_STATE_POWEROFF MAV_STATE = 7
	// System is terminating itself (failsafe or commanded).
	MAV_STATE_FLIGHT_TERMINATION MAV_STATE = 8
)

var labels_MAV_STATE = map[MAV_STATE]string{
	MAV_STATE_UNINIT:             "MAV_STATE_UNINIT",
	MAV_STATE_BOOT:               "MAV_STATE_BOOT",
	MAV_STATE_CALIBRATING:        "MAV_STATE_CALIBRATING",
	MAV_STATE_STANDBY:            "MAV_STATE_STANDBY",
	MAV_STATE_ACTIVE:             "MAV_STATE_ACTIVE",
	MAV_STATE_CRITICAL:           "MAV_STATE_CRITICAL",
	MAV_STATE_EMERGENCY:          "MAV_STATE_EMERGENCY",
	MAV_STATE_POWEROFF:           "MAV_STATE_POWEROFF",
	MAV_STATE_FLIGHT_TERMINATION: "MAV_STATE_FLIGHT_TERMINATION",
}

var values_MAV_STATE = map[string]MAV_STATE{
	"MAV_STATE_UNINIT":             MAV_STATE_UNINIT,
	"MAV_STATE_BOOT":               MAV_STATE_BOOT,
	"MAV_STATE_CALIBRATING":        MAV_STATE_CALIBRATING,
	"MAV_STATE_STANDBY":            MAV_STATE_STANDBY,
	"MAV_STATE_ACTIVE":             MAV_STATE_ACTIVE,
	"MAV_STATE_CRITICAL":           MAV_STATE_CRITICAL,
	"MAV_STATE_EMERGENCY":          MAV_STATE_EMERGENCY,
	"MAV_STATE_POWEROFF":           MAV_STATE_POWEROFF,
	"MAV_STATE_FLIGHT_TERMINATION": MAV_STATE_FLIGHT_TERMINATION,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_STATE) MarshalText() ([]byte, error) {
	name, ok := labels_MAV_STATE[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_STATE) UnmarshalText(text []byte) error {
	value, ok := values_MAV_STATE[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_STATE) String() string {
	name, ok := labels_MAV_STATE[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
