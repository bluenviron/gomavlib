//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl
package matrixpilot

import (
	"errors"
)

type GPS_INPUT_IGNORE_FLAGS uint32

const (
	// ignore altitude field
	GPS_INPUT_IGNORE_FLAG_ALT GPS_INPUT_IGNORE_FLAGS = 1
	// ignore hdop field
	GPS_INPUT_IGNORE_FLAG_HDOP GPS_INPUT_IGNORE_FLAGS = 2
	// ignore vdop field
	GPS_INPUT_IGNORE_FLAG_VDOP GPS_INPUT_IGNORE_FLAGS = 4
	// ignore horizontal velocity field (vn and ve)
	GPS_INPUT_IGNORE_FLAG_VEL_HORIZ GPS_INPUT_IGNORE_FLAGS = 8
	// ignore vertical velocity field (vd)
	GPS_INPUT_IGNORE_FLAG_VEL_VERT GPS_INPUT_IGNORE_FLAGS = 16
	// ignore speed accuracy field
	GPS_INPUT_IGNORE_FLAG_SPEED_ACCURACY GPS_INPUT_IGNORE_FLAGS = 32
	// ignore horizontal accuracy field
	GPS_INPUT_IGNORE_FLAG_HORIZONTAL_ACCURACY GPS_INPUT_IGNORE_FLAGS = 64
	// ignore vertical accuracy field
	GPS_INPUT_IGNORE_FLAG_VERTICAL_ACCURACY GPS_INPUT_IGNORE_FLAGS = 128
)

var labels_GPS_INPUT_IGNORE_FLAGS = map[GPS_INPUT_IGNORE_FLAGS]string{}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GPS_INPUT_IGNORE_FLAGS) MarshalText() ([]byte, error) {
	if l, ok := labels_GPS_INPUT_IGNORE_FLAGS[e]; ok {
		return []byte(l), nil
	}
	return nil, errors.New("invalid value")
}

var reverseLabels_GPS_INPUT_IGNORE_FLAGS = map[string]GPS_INPUT_IGNORE_FLAGS{}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GPS_INPUT_IGNORE_FLAGS) UnmarshalText(text []byte) error {
	if rl, ok := reverseLabels_GPS_INPUT_IGNORE_FLAGS[string(text)]; ok {
		*e = rl
		return nil
	}
	return errors.New("invalid value")
}

// String implements the fmt.Stringer interface.
func (e GPS_INPUT_IGNORE_FLAGS) String() string {
	if l, ok := labels_GPS_INPUT_IGNORE_FLAGS[e]; ok {
		return l
	}
	return "invalid value"
}