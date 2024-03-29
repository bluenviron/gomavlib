//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package csairlink

import (
	"fmt"
	"strconv"
)

type AIRLINK_EYE_TURN_INIT_TYPE uint64

const (
	AIRLINK_TURN_INIT_START AIRLINK_EYE_TURN_INIT_TYPE = 0
	AIRLINK_TURN_INIT_OK    AIRLINK_EYE_TURN_INIT_TYPE = 1
	AIRLINK_TURN_INIT_BAD   AIRLINK_EYE_TURN_INIT_TYPE = 2
)

var labels_AIRLINK_EYE_TURN_INIT_TYPE = map[AIRLINK_EYE_TURN_INIT_TYPE]string{
	AIRLINK_TURN_INIT_START: "AIRLINK_TURN_INIT_START",
	AIRLINK_TURN_INIT_OK:    "AIRLINK_TURN_INIT_OK",
	AIRLINK_TURN_INIT_BAD:   "AIRLINK_TURN_INIT_BAD",
}

var values_AIRLINK_EYE_TURN_INIT_TYPE = map[string]AIRLINK_EYE_TURN_INIT_TYPE{
	"AIRLINK_TURN_INIT_START": AIRLINK_TURN_INIT_START,
	"AIRLINK_TURN_INIT_OK":    AIRLINK_TURN_INIT_OK,
	"AIRLINK_TURN_INIT_BAD":   AIRLINK_TURN_INIT_BAD,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e AIRLINK_EYE_TURN_INIT_TYPE) MarshalText() ([]byte, error) {
	if name, ok := labels_AIRLINK_EYE_TURN_INIT_TYPE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *AIRLINK_EYE_TURN_INIT_TYPE) UnmarshalText(text []byte) error {
	if value, ok := values_AIRLINK_EYE_TURN_INIT_TYPE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = AIRLINK_EYE_TURN_INIT_TYPE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e AIRLINK_EYE_TURN_INIT_TYPE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
