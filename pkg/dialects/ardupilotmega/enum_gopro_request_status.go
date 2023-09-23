//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
)

type GOPRO_REQUEST_STATUS uint32

const (
	// The write message with ID indicated succeeded.
	GOPRO_REQUEST_SUCCESS GOPRO_REQUEST_STATUS = 0
	// The write message with ID indicated failed.
	GOPRO_REQUEST_FAILED GOPRO_REQUEST_STATUS = 1
)

var labels_GOPRO_REQUEST_STATUS = map[GOPRO_REQUEST_STATUS]string{
	GOPRO_REQUEST_SUCCESS: "GOPRO_REQUEST_SUCCESS",
	GOPRO_REQUEST_FAILED:  "GOPRO_REQUEST_FAILED",
}

var values_GOPRO_REQUEST_STATUS = map[string]GOPRO_REQUEST_STATUS{
	"GOPRO_REQUEST_SUCCESS": GOPRO_REQUEST_SUCCESS,
	"GOPRO_REQUEST_FAILED":  GOPRO_REQUEST_FAILED,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GOPRO_REQUEST_STATUS) MarshalText() ([]byte, error) {
	name, ok := labels_GOPRO_REQUEST_STATUS[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GOPRO_REQUEST_STATUS) UnmarshalText(text []byte) error {
	value, ok := values_GOPRO_REQUEST_STATUS[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e GOPRO_REQUEST_STATUS) String() string {
	name, ok := labels_GOPRO_REQUEST_STATUS[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
