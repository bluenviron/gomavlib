//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package csairlink

import (
	"fmt"
	"strconv"
)

type AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE uint32

const (
	AIRLINK_HPR_PARTNER_NOT_READY AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE = 0
	AIRLINK_HPR_PARTNER_READY     AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE = 1
)

var labels_AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE = map[AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE]string{
	AIRLINK_HPR_PARTNER_NOT_READY: "AIRLINK_HPR_PARTNER_NOT_READY",
	AIRLINK_HPR_PARTNER_READY:     "AIRLINK_HPR_PARTNER_READY",
}

var values_AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE = map[string]AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE{
	"AIRLINK_HPR_PARTNER_NOT_READY": AIRLINK_HPR_PARTNER_NOT_READY,
	"AIRLINK_HPR_PARTNER_READY":     AIRLINK_HPR_PARTNER_READY,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE) MarshalText() ([]byte, error) {
	if name, ok := labels_AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE) UnmarshalText(text []byte) error {
	if value, ok := values_AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e AIRLINK_EYE_GS_HOLE_PUSH_RESP_TYPE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}