//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
	"strings"
)

// These encode the sensors whose status is sent as part of the SYS_STATUS message in the extended fields.
type MAV_SYS_STATUS_SENSOR_EXTENDED uint32

const (
	// 0x01 Recovery system (parachute, balloon, retracts etc)
	MAV_SYS_STATUS_RECOVERY_SYSTEM MAV_SYS_STATUS_SENSOR_EXTENDED = 1
)

var labels_MAV_SYS_STATUS_SENSOR_EXTENDED = map[MAV_SYS_STATUS_SENSOR_EXTENDED]string{
	MAV_SYS_STATUS_RECOVERY_SYSTEM: "MAV_SYS_STATUS_RECOVERY_SYSTEM",
}

var values_MAV_SYS_STATUS_SENSOR_EXTENDED = map[string]MAV_SYS_STATUS_SENSOR_EXTENDED{
	"MAV_SYS_STATUS_RECOVERY_SYSTEM": MAV_SYS_STATUS_RECOVERY_SYSTEM,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e MAV_SYS_STATUS_SENSOR_EXTENDED) MarshalText() ([]byte, error) {
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for i := 0; i < 1; i++ {
		mask := MAV_SYS_STATUS_SENSOR_EXTENDED(1 << i)
		if e&mask == mask {
			names = append(names, labels_MAV_SYS_STATUS_SENSOR_EXTENDED[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *MAV_SYS_STATUS_SENSOR_EXTENDED) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask MAV_SYS_STATUS_SENSOR_EXTENDED
	for _, label := range labels {
		if value, ok := values_MAV_SYS_STATUS_SENSOR_EXTENDED[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= MAV_SYS_STATUS_SENSOR_EXTENDED(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e MAV_SYS_STATUS_SENSOR_EXTENDED) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
