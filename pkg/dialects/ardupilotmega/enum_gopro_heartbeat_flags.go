//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"fmt"
	"strconv"
	"strings"
)

type GOPRO_HEARTBEAT_FLAGS uint32

const (
	// GoPro is currently recording.
	GOPRO_FLAG_RECORDING GOPRO_HEARTBEAT_FLAGS = 1
)

var labels_GOPRO_HEARTBEAT_FLAGS = map[GOPRO_HEARTBEAT_FLAGS]string{
	GOPRO_FLAG_RECORDING: "GOPRO_FLAG_RECORDING",
}

var values_GOPRO_HEARTBEAT_FLAGS = map[string]GOPRO_HEARTBEAT_FLAGS{
	"GOPRO_FLAG_RECORDING": GOPRO_FLAG_RECORDING,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e GOPRO_HEARTBEAT_FLAGS) MarshalText() ([]byte, error) {
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for i := 0; i < 1; i++ {
		mask := GOPRO_HEARTBEAT_FLAGS(1 << i)
		if e&mask == mask {
			names = append(names, labels_GOPRO_HEARTBEAT_FLAGS[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *GOPRO_HEARTBEAT_FLAGS) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask GOPRO_HEARTBEAT_FLAGS
	for _, label := range labels {
		if value, ok := values_GOPRO_HEARTBEAT_FLAGS[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= GOPRO_HEARTBEAT_FLAGS(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	*e = mask
	return nil
}

// String implements the fmt.Stringer interface.
func (e GOPRO_HEARTBEAT_FLAGS) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
