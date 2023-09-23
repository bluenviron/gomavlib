//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
)

type HEADING_TYPE uint32

const (
	HEADING_TYPE_COURSE_OVER_GROUND HEADING_TYPE = 0
	HEADING_TYPE_HEADING            HEADING_TYPE = 1
)

var labels_HEADING_TYPE = map[HEADING_TYPE]string{
	HEADING_TYPE_COURSE_OVER_GROUND: "HEADING_TYPE_COURSE_OVER_GROUND",
	HEADING_TYPE_HEADING:            "HEADING_TYPE_HEADING",
}

var values_HEADING_TYPE = map[string]HEADING_TYPE{
	"HEADING_TYPE_COURSE_OVER_GROUND": HEADING_TYPE_COURSE_OVER_GROUND,
	"HEADING_TYPE_HEADING":            HEADING_TYPE_HEADING,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e HEADING_TYPE) MarshalText() ([]byte, error) {
	name, ok := labels_HEADING_TYPE[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *HEADING_TYPE) UnmarshalText(text []byte) error {
	value, ok := values_HEADING_TYPE[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e HEADING_TYPE) String() string {
	name, ok := labels_HEADING_TYPE[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
