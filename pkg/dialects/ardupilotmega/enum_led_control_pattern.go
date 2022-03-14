//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl
package ardupilotmega

import (
	"errors"
)

type LED_CONTROL_PATTERN uint32

const (
	// LED patterns off (return control to regular vehicle control).
	LED_CONTROL_PATTERN_OFF LED_CONTROL_PATTERN = 0
	// LEDs show pattern during firmware update.
	LED_CONTROL_PATTERN_FIRMWAREUPDATE LED_CONTROL_PATTERN = 1
	// Custom Pattern using custom bytes fields.
	LED_CONTROL_PATTERN_CUSTOM LED_CONTROL_PATTERN = 255
)

var labels_LED_CONTROL_PATTERN = map[LED_CONTROL_PATTERN]string{}

// MarshalText implements the encoding.TextMarshaler interface.
func (e LED_CONTROL_PATTERN) MarshalText() ([]byte, error) {
	if l, ok := labels_LED_CONTROL_PATTERN[e]; ok {
		return []byte(l), nil
	}
	return nil, errors.New("invalid value")
}

var reverseLabels_LED_CONTROL_PATTERN = map[string]LED_CONTROL_PATTERN{}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *LED_CONTROL_PATTERN) UnmarshalText(text []byte) error {
	if rl, ok := reverseLabels_LED_CONTROL_PATTERN[string(text)]; ok {
		*e = rl
		return nil
	}
	return errors.New("invalid value")
}

// String implements the fmt.Stringer interface.
func (e LED_CONTROL_PATTERN) String() string {
	if l, ok := labels_LED_CONTROL_PATTERN[e]; ok {
		return l
	}
	return "invalid value"
}