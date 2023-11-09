//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
)

type GOPRO_PHOTO_RESOLUTION uint32

const (
	// 5MP Medium.
	GOPRO_PHOTO_RESOLUTION_5MP_MEDIUM GOPRO_PHOTO_RESOLUTION = 0
	// 7MP Medium.
	GOPRO_PHOTO_RESOLUTION_7MP_MEDIUM GOPRO_PHOTO_RESOLUTION = 1
	// 7MP Wide.
	GOPRO_PHOTO_RESOLUTION_7MP_WIDE GOPRO_PHOTO_RESOLUTION = 2
	// 10MP Wide.
	GOPRO_PHOTO_RESOLUTION_10MP_WIDE GOPRO_PHOTO_RESOLUTION = 3
	// 12MP Wide.
	GOPRO_PHOTO_RESOLUTION_12MP_WIDE GOPRO_PHOTO_RESOLUTION = 4
)

var labels_GOPRO_PHOTO_RESOLUTION = map[GOPRO_PHOTO_RESOLUTION]string{
	GOPRO_PHOTO_RESOLUTION_5MP_MEDIUM: "GOPRO_PHOTO_RESOLUTION_5MP_MEDIUM",
	GOPRO_PHOTO_RESOLUTION_7MP_MEDIUM: "GOPRO_PHOTO_RESOLUTION_7MP_MEDIUM",
	GOPRO_PHOTO_RESOLUTION_7MP_WIDE:   "GOPRO_PHOTO_RESOLUTION_7MP_WIDE",
	GOPRO_PHOTO_RESOLUTION_10MP_WIDE:  "GOPRO_PHOTO_RESOLUTION_10MP_WIDE",
	GOPRO_PHOTO_RESOLUTION_12MP_WIDE:  "GOPRO_PHOTO_RESOLUTION_12MP_WIDE",
}

var values_GOPRO_PHOTO_RESOLUTION = map[string]GOPRO_PHOTO_RESOLUTION{
	"GOPRO_PHOTO_RESOLUTION_5MP_MEDIUM": GOPRO_PHOTO_RESOLUTION_5MP_MEDIUM,
	"GOPRO_PHOTO_RESOLUTION_7MP_MEDIUM": GOPRO_PHOTO_RESOLUTION_7MP_MEDIUM,
	"GOPRO_PHOTO_RESOLUTION_7MP_WIDE":   GOPRO_PHOTO_RESOLUTION_7MP_WIDE,
	"GOPRO_PHOTO_RESOLUTION_10MP_WIDE":  GOPRO_PHOTO_RESOLUTION_10MP_WIDE,
	"GOPRO_PHOTO_RESOLUTION_12MP_WIDE":  GOPRO_PHOTO_RESOLUTION_12MP_WIDE,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GOPRO_PHOTO_RESOLUTION) MarshalText() ([]byte, error) {
	if name, ok := labels_GOPRO_PHOTO_RESOLUTION[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GOPRO_PHOTO_RESOLUTION) UnmarshalText(text []byte) error {
	if value, ok := values_GOPRO_PHOTO_RESOLUTION[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = GOPRO_PHOTO_RESOLUTION(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e GOPRO_PHOTO_RESOLUTION) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
