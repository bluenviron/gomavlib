//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
	"strings"
)

// Bitmap of options for the MAV_CMD_DO_REPOSITION
type MAV_DO_REPOSITION_FLAGS uint64

const (
	// The aircraft should immediately transition into guided. This should not be set for follow me applications
	MAV_DO_REPOSITION_FLAGS_CHANGE_MODE MAV_DO_REPOSITION_FLAGS = 1
	// Yaw relative to the vehicle current heading (if not set, relative to North).
	MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW MAV_DO_REPOSITION_FLAGS = 2
)

var values_MAV_DO_REPOSITION_FLAGS = []MAV_DO_REPOSITION_FLAGS{
	MAV_DO_REPOSITION_FLAGS_CHANGE_MODE,
	MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW,
}

var value_to_label_MAV_DO_REPOSITION_FLAGS = map[MAV_DO_REPOSITION_FLAGS]string{
	MAV_DO_REPOSITION_FLAGS_CHANGE_MODE:  "MAV_DO_REPOSITION_FLAGS_CHANGE_MODE",
	MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW: "MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW",
}

var label_to_value_MAV_DO_REPOSITION_FLAGS = map[string]MAV_DO_REPOSITION_FLAGS{
	"MAV_DO_REPOSITION_FLAGS_CHANGE_MODE":  MAV_DO_REPOSITION_FLAGS_CHANGE_MODE,
	"MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW": MAV_DO_REPOSITION_FLAGS_RELATIVE_YAW,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_DO_REPOSITION_FLAGS) MarshalText() ([]byte, error) {
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for _, val := range values_MAV_DO_REPOSITION_FLAGS {
		if e&val == val {
			names = append(names, value_to_label_MAV_DO_REPOSITION_FLAGS[val])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_DO_REPOSITION_FLAGS) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask MAV_DO_REPOSITION_FLAGS
	for _, label := range labels {
		if value, ok := label_to_value_MAV_DO_REPOSITION_FLAGS[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= MAV_DO_REPOSITION_FLAGS(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	*e = mask
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_DO_REPOSITION_FLAGS) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
