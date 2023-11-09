//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package asluav

import (
	"fmt"
	"strconv"
)

type GSM_LINK_TYPE uint32

const (
	// no service
	GSM_LINK_TYPE_NONE GSM_LINK_TYPE = 0
	// link type unknown
	GSM_LINK_TYPE_UNKNOWN GSM_LINK_TYPE = 1
	// 2G (GSM/GRPS/EDGE) link
	GSM_LINK_TYPE_2G GSM_LINK_TYPE = 2
	// 3G link (WCDMA/HSDPA/HSPA)
	GSM_LINK_TYPE_3G GSM_LINK_TYPE = 3
	// 4G link (LTE)
	GSM_LINK_TYPE_4G GSM_LINK_TYPE = 4
)

var labels_GSM_LINK_TYPE = map[GSM_LINK_TYPE]string{
	GSM_LINK_TYPE_NONE:    "GSM_LINK_TYPE_NONE",
	GSM_LINK_TYPE_UNKNOWN: "GSM_LINK_TYPE_UNKNOWN",
	GSM_LINK_TYPE_2G:      "GSM_LINK_TYPE_2G",
	GSM_LINK_TYPE_3G:      "GSM_LINK_TYPE_3G",
	GSM_LINK_TYPE_4G:      "GSM_LINK_TYPE_4G",
}

var values_GSM_LINK_TYPE = map[string]GSM_LINK_TYPE{
	"GSM_LINK_TYPE_NONE":    GSM_LINK_TYPE_NONE,
	"GSM_LINK_TYPE_UNKNOWN": GSM_LINK_TYPE_UNKNOWN,
	"GSM_LINK_TYPE_2G":      GSM_LINK_TYPE_2G,
	"GSM_LINK_TYPE_3G":      GSM_LINK_TYPE_3G,
	"GSM_LINK_TYPE_4G":      GSM_LINK_TYPE_4G,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GSM_LINK_TYPE) MarshalText() ([]byte, error) {
	if name, ok := labels_GSM_LINK_TYPE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GSM_LINK_TYPE) UnmarshalText(text []byte) error {
	if value, ok := values_GSM_LINK_TYPE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = GSM_LINK_TYPE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e GSM_LINK_TYPE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
