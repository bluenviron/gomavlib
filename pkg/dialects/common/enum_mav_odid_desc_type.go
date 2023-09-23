//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

type MAV_ODID_DESC_TYPE uint32

const (
	// Optional free-form text description of the purpose of the flight.
	MAV_ODID_DESC_TYPE_TEXT MAV_ODID_DESC_TYPE = 0
	// Optional additional clarification when status == MAV_ODID_STATUS_EMERGENCY.
	MAV_ODID_DESC_TYPE_EMERGENCY MAV_ODID_DESC_TYPE = 1
	// Optional additional clarification when status != MAV_ODID_STATUS_EMERGENCY.
	MAV_ODID_DESC_TYPE_EXTENDED_STATUS MAV_ODID_DESC_TYPE = 2
)

var labels_MAV_ODID_DESC_TYPE = map[MAV_ODID_DESC_TYPE]string{
	MAV_ODID_DESC_TYPE_TEXT:            "MAV_ODID_DESC_TYPE_TEXT",
	MAV_ODID_DESC_TYPE_EMERGENCY:       "MAV_ODID_DESC_TYPE_EMERGENCY",
	MAV_ODID_DESC_TYPE_EXTENDED_STATUS: "MAV_ODID_DESC_TYPE_EXTENDED_STATUS",
}

var values_MAV_ODID_DESC_TYPE = map[string]MAV_ODID_DESC_TYPE{
	"MAV_ODID_DESC_TYPE_TEXT":            MAV_ODID_DESC_TYPE_TEXT,
	"MAV_ODID_DESC_TYPE_EMERGENCY":       MAV_ODID_DESC_TYPE_EMERGENCY,
	"MAV_ODID_DESC_TYPE_EXTENDED_STATUS": MAV_ODID_DESC_TYPE_EXTENDED_STATUS,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_ODID_DESC_TYPE) MarshalText() ([]byte, error) {
	name, ok := labels_MAV_ODID_DESC_TYPE[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_ODID_DESC_TYPE) UnmarshalText(text []byte) error {
	value, ok := values_MAV_ODID_DESC_TYPE[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_ODID_DESC_TYPE) String() string {
	name, ok := labels_MAV_ODID_DESC_TYPE[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
