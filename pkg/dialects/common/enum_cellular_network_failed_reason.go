//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// These flags are used to diagnose the failure state of CELLULAR_STATUS
type CELLULAR_NETWORK_FAILED_REASON uint32

const (
	// No error
	CELLULAR_NETWORK_FAILED_REASON_NONE CELLULAR_NETWORK_FAILED_REASON = 0
	// Error state is unknown
	CELLULAR_NETWORK_FAILED_REASON_UNKNOWN CELLULAR_NETWORK_FAILED_REASON = 1
	// SIM is required for the modem but missing
	CELLULAR_NETWORK_FAILED_REASON_SIM_MISSING CELLULAR_NETWORK_FAILED_REASON = 2
	// SIM is available, but not usable for connection
	CELLULAR_NETWORK_FAILED_REASON_SIM_ERROR CELLULAR_NETWORK_FAILED_REASON = 3
)

var labels_CELLULAR_NETWORK_FAILED_REASON = map[CELLULAR_NETWORK_FAILED_REASON]string{
	CELLULAR_NETWORK_FAILED_REASON_NONE:        "CELLULAR_NETWORK_FAILED_REASON_NONE",
	CELLULAR_NETWORK_FAILED_REASON_UNKNOWN:     "CELLULAR_NETWORK_FAILED_REASON_UNKNOWN",
	CELLULAR_NETWORK_FAILED_REASON_SIM_MISSING: "CELLULAR_NETWORK_FAILED_REASON_SIM_MISSING",
	CELLULAR_NETWORK_FAILED_REASON_SIM_ERROR:   "CELLULAR_NETWORK_FAILED_REASON_SIM_ERROR",
}

var values_CELLULAR_NETWORK_FAILED_REASON = map[string]CELLULAR_NETWORK_FAILED_REASON{
	"CELLULAR_NETWORK_FAILED_REASON_NONE":        CELLULAR_NETWORK_FAILED_REASON_NONE,
	"CELLULAR_NETWORK_FAILED_REASON_UNKNOWN":     CELLULAR_NETWORK_FAILED_REASON_UNKNOWN,
	"CELLULAR_NETWORK_FAILED_REASON_SIM_MISSING": CELLULAR_NETWORK_FAILED_REASON_SIM_MISSING,
	"CELLULAR_NETWORK_FAILED_REASON_SIM_ERROR":   CELLULAR_NETWORK_FAILED_REASON_SIM_ERROR,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e CELLULAR_NETWORK_FAILED_REASON) MarshalText() ([]byte, error) {
	name, ok := labels_CELLULAR_NETWORK_FAILED_REASON[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *CELLULAR_NETWORK_FAILED_REASON) UnmarshalText(text []byte) error {
	value, ok := values_CELLULAR_NETWORK_FAILED_REASON[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e CELLULAR_NETWORK_FAILED_REASON) String() string {
	name, ok := labels_CELLULAR_NETWORK_FAILED_REASON[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
