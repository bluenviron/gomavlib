//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
)

type GOPRO_PROTUNE_COLOUR uint32

const (
	// Auto.
	GOPRO_PROTUNE_COLOUR_STANDARD GOPRO_PROTUNE_COLOUR = 0
	// Neutral.
	GOPRO_PROTUNE_COLOUR_NEUTRAL GOPRO_PROTUNE_COLOUR = 1
)

var labels_GOPRO_PROTUNE_COLOUR = map[GOPRO_PROTUNE_COLOUR]string{
	GOPRO_PROTUNE_COLOUR_STANDARD: "GOPRO_PROTUNE_COLOUR_STANDARD",
	GOPRO_PROTUNE_COLOUR_NEUTRAL:  "GOPRO_PROTUNE_COLOUR_NEUTRAL",
}

var values_GOPRO_PROTUNE_COLOUR = map[string]GOPRO_PROTUNE_COLOUR{
	"GOPRO_PROTUNE_COLOUR_STANDARD": GOPRO_PROTUNE_COLOUR_STANDARD,
	"GOPRO_PROTUNE_COLOUR_NEUTRAL":  GOPRO_PROTUNE_COLOUR_NEUTRAL,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GOPRO_PROTUNE_COLOUR) MarshalText() ([]byte, error) {
	name, ok := labels_GOPRO_PROTUNE_COLOUR[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GOPRO_PROTUNE_COLOUR) UnmarshalText(text []byte) error {
	value, ok := values_GOPRO_PROTUNE_COLOUR[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e GOPRO_PROTUNE_COLOUR) String() string {
	name, ok := labels_GOPRO_PROTUNE_COLOUR[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
