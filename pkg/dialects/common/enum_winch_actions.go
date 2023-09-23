//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Winch actions.
type WINCH_ACTIONS uint32

const (
	// Allow motor to freewheel.
	WINCH_RELAXED WINCH_ACTIONS = 0
	// Wind or unwind specified length of line, optionally using specified rate.
	WINCH_RELATIVE_LENGTH_CONTROL WINCH_ACTIONS = 1
	// Wind or unwind line at specified rate.
	WINCH_RATE_CONTROL WINCH_ACTIONS = 2
	// Perform the locking sequence to relieve motor while in the fully retracted position. Only action and instance command parameters are used, others are ignored.
	WINCH_LOCK WINCH_ACTIONS = 3
	// Sequence of drop, slow down, touch down, reel up, lock. Only action and instance command parameters are used, others are ignored.
	WINCH_DELIVER WINCH_ACTIONS = 4
	// Engage motor and hold current position. Only action and instance command parameters are used, others are ignored.
	WINCH_HOLD WINCH_ACTIONS = 5
	// Return the reel to the fully retracted position. Only action and instance command parameters are used, others are ignored.
	WINCH_RETRACT WINCH_ACTIONS = 6
	// Load the reel with line. The winch will calculate the total loaded length and stop when the tension exceeds a threshold. Only action and instance command parameters are used, others are ignored.
	WINCH_LOAD_LINE WINCH_ACTIONS = 7
	// Spool out the entire length of the line. Only action and instance command parameters are used, others are ignored.
	WINCH_ABANDON_LINE WINCH_ACTIONS = 8
	// Spools out just enough to present the hook to the user to load the payload. Only action and instance command parameters are used, others are ignored
	WINCH_LOAD_PAYLOAD WINCH_ACTIONS = 9
)

var labels_WINCH_ACTIONS = map[WINCH_ACTIONS]string{
	WINCH_RELAXED:                 "WINCH_RELAXED",
	WINCH_RELATIVE_LENGTH_CONTROL: "WINCH_RELATIVE_LENGTH_CONTROL",
	WINCH_RATE_CONTROL:            "WINCH_RATE_CONTROL",
	WINCH_LOCK:                    "WINCH_LOCK",
	WINCH_DELIVER:                 "WINCH_DELIVER",
	WINCH_HOLD:                    "WINCH_HOLD",
	WINCH_RETRACT:                 "WINCH_RETRACT",
	WINCH_LOAD_LINE:               "WINCH_LOAD_LINE",
	WINCH_ABANDON_LINE:            "WINCH_ABANDON_LINE",
	WINCH_LOAD_PAYLOAD:            "WINCH_LOAD_PAYLOAD",
}

var values_WINCH_ACTIONS = map[string]WINCH_ACTIONS{
	"WINCH_RELAXED":                 WINCH_RELAXED,
	"WINCH_RELATIVE_LENGTH_CONTROL": WINCH_RELATIVE_LENGTH_CONTROL,
	"WINCH_RATE_CONTROL":            WINCH_RATE_CONTROL,
	"WINCH_LOCK":                    WINCH_LOCK,
	"WINCH_DELIVER":                 WINCH_DELIVER,
	"WINCH_HOLD":                    WINCH_HOLD,
	"WINCH_RETRACT":                 WINCH_RETRACT,
	"WINCH_LOAD_LINE":               WINCH_LOAD_LINE,
	"WINCH_ABANDON_LINE":            WINCH_ABANDON_LINE,
	"WINCH_LOAD_PAYLOAD":            WINCH_LOAD_PAYLOAD,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e WINCH_ACTIONS) MarshalText() ([]byte, error) {
	name, ok := labels_WINCH_ACTIONS[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *WINCH_ACTIONS) UnmarshalText(text []byte) error {
	value, ok := values_WINCH_ACTIONS[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e WINCH_ACTIONS) String() string {
	name, ok := labels_WINCH_ACTIONS[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
