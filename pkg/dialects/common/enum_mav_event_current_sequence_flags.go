//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Flags for CURRENT_EVENT_SEQUENCE.
type MAV_EVENT_CURRENT_SEQUENCE_FLAGS uint32

const (
	// A sequence reset has happened (e.g. vehicle reboot).
	MAV_EVENT_CURRENT_SEQUENCE_FLAGS_RESET MAV_EVENT_CURRENT_SEQUENCE_FLAGS = 1
)

var labels_MAV_EVENT_CURRENT_SEQUENCE_FLAGS = map[MAV_EVENT_CURRENT_SEQUENCE_FLAGS]string{
	MAV_EVENT_CURRENT_SEQUENCE_FLAGS_RESET: "MAV_EVENT_CURRENT_SEQUENCE_FLAGS_RESET",
}

var values_MAV_EVENT_CURRENT_SEQUENCE_FLAGS = map[string]MAV_EVENT_CURRENT_SEQUENCE_FLAGS{
	"MAV_EVENT_CURRENT_SEQUENCE_FLAGS_RESET": MAV_EVENT_CURRENT_SEQUENCE_FLAGS_RESET,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_EVENT_CURRENT_SEQUENCE_FLAGS) MarshalText() ([]byte, error) {
	name, ok := labels_MAV_EVENT_CURRENT_SEQUENCE_FLAGS[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_EVENT_CURRENT_SEQUENCE_FLAGS) UnmarshalText(text []byte) error {
	value, ok := values_MAV_EVENT_CURRENT_SEQUENCE_FLAGS[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_EVENT_CURRENT_SEQUENCE_FLAGS) String() string {
	name, ok := labels_MAV_EVENT_CURRENT_SEQUENCE_FLAGS[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
