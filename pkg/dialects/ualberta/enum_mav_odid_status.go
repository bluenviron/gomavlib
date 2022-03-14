//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl
package ualberta

import (
	"errors"
)

type MAV_ODID_STATUS uint32

const (
	// The status of the (UA) Unmanned Aircraft is undefined.
	MAV_ODID_STATUS_UNDECLARED MAV_ODID_STATUS = 0
	// The UA is on the ground.
	MAV_ODID_STATUS_GROUND MAV_ODID_STATUS = 1
	// The UA is in the air.
	MAV_ODID_STATUS_AIRBORNE MAV_ODID_STATUS = 2
	// The UA is having an emergency.
	MAV_ODID_STATUS_EMERGENCY MAV_ODID_STATUS = 3
)

var labels_MAV_ODID_STATUS = map[MAV_ODID_STATUS]string{}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_ODID_STATUS) MarshalText() ([]byte, error) {
	if l, ok := labels_MAV_ODID_STATUS[e]; ok {
		return []byte(l), nil
	}
	return nil, errors.New("invalid value")
}

var reverseLabels_MAV_ODID_STATUS = map[string]MAV_ODID_STATUS{}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_ODID_STATUS) UnmarshalText(text []byte) error {
	if rl, ok := reverseLabels_MAV_ODID_STATUS[string(text)]; ok {
		*e = rl
		return nil
	}
	return errors.New("invalid value")
}

// String implements the fmt.Stringer interface.
func (e MAV_ODID_STATUS) String() string {
	if l, ok := labels_MAV_ODID_STATUS[e]; ok {
		return l
	}
	return "invalid value"
}