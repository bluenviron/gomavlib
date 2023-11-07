//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Actions being taken to mitigate/prevent fence breach
type FENCE_MITIGATE uint32

const (
	// Unknown
	FENCE_MITIGATE_UNKNOWN FENCE_MITIGATE = 0
	// No actions being taken
	FENCE_MITIGATE_NONE FENCE_MITIGATE = 1
	// Velocity limiting active to prevent breach
	FENCE_MITIGATE_VEL_LIMIT FENCE_MITIGATE = 2
)

var labels_FENCE_MITIGATE = map[FENCE_MITIGATE]string{
	FENCE_MITIGATE_UNKNOWN:   "FENCE_MITIGATE_UNKNOWN",
	FENCE_MITIGATE_NONE:      "FENCE_MITIGATE_NONE",
	FENCE_MITIGATE_VEL_LIMIT: "FENCE_MITIGATE_VEL_LIMIT",
}

var values_FENCE_MITIGATE = map[string]FENCE_MITIGATE{
	"FENCE_MITIGATE_UNKNOWN":   FENCE_MITIGATE_UNKNOWN,
	"FENCE_MITIGATE_NONE":      FENCE_MITIGATE_NONE,
	"FENCE_MITIGATE_VEL_LIMIT": FENCE_MITIGATE_VEL_LIMIT,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e FENCE_MITIGATE) MarshalText() ([]byte, error) {
	if name, ok := labels_FENCE_MITIGATE[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *FENCE_MITIGATE) UnmarshalText(text []byte) error {
	if value, ok := values_FENCE_MITIGATE[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = FENCE_MITIGATE(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e FENCE_MITIGATE) String() string {
	if name, ok := labels_FENCE_MITIGATE[e]; ok {
		return name
	}
	return strconv.Itoa(int(e))
}
