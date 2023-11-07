//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
	"strings"
)

type GOPRO_VIDEO_SETTINGS_FLAGS uint32

const (
	// 0=NTSC, 1=PAL.
	GOPRO_VIDEO_SETTINGS_TV_MODE GOPRO_VIDEO_SETTINGS_FLAGS = 1
)

var labels_GOPRO_VIDEO_SETTINGS_FLAGS = map[GOPRO_VIDEO_SETTINGS_FLAGS]string{
	GOPRO_VIDEO_SETTINGS_TV_MODE: "GOPRO_VIDEO_SETTINGS_TV_MODE",
}

var values_GOPRO_VIDEO_SETTINGS_FLAGS = map[string]GOPRO_VIDEO_SETTINGS_FLAGS{
	"GOPRO_VIDEO_SETTINGS_TV_MODE": GOPRO_VIDEO_SETTINGS_TV_MODE,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GOPRO_VIDEO_SETTINGS_FLAGS) MarshalText() ([]byte, error) {
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for i := 0; i < 1; i++ {
		mask := GOPRO_VIDEO_SETTINGS_FLAGS(1 << i)
		if e&mask == mask {
			names = append(names, labels_GOPRO_VIDEO_SETTINGS_FLAGS[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GOPRO_VIDEO_SETTINGS_FLAGS) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask GOPRO_VIDEO_SETTINGS_FLAGS
	for _, label := range labels {
		if value, ok := values_GOPRO_VIDEO_SETTINGS_FLAGS[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= GOPRO_VIDEO_SETTINGS_FLAGS(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e GOPRO_VIDEO_SETTINGS_FLAGS) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
