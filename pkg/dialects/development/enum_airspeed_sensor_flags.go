//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package development

import (
	"fmt"
	"strings"
)

// Airspeed sensor flags
type AIRSPEED_SENSOR_FLAGS uint32

const (
	// Airspeed sensor is unhealthy
	AIRSPEED_SENSOR_UNHEALTHY AIRSPEED_SENSOR_FLAGS = 0
	// True if the data from this sensor is being actively used by the flight controller for guidance, navigation or control.
	AIRSPEED_SENSOR_USING AIRSPEED_SENSOR_FLAGS = 1
)

var labels_AIRSPEED_SENSOR_FLAGS = map[AIRSPEED_SENSOR_FLAGS]string{
	AIRSPEED_SENSOR_UNHEALTHY: "AIRSPEED_SENSOR_UNHEALTHY",
	AIRSPEED_SENSOR_USING:     "AIRSPEED_SENSOR_USING",
}

var values_AIRSPEED_SENSOR_FLAGS = map[string]AIRSPEED_SENSOR_FLAGS{
	"AIRSPEED_SENSOR_UNHEALTHY": AIRSPEED_SENSOR_UNHEALTHY,
	"AIRSPEED_SENSOR_USING":     AIRSPEED_SENSOR_USING,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e AIRSPEED_SENSOR_FLAGS) MarshalText() ([]byte, error) {
	var names []string
	for i := 0; i < 2; i++ {
		mask := AIRSPEED_SENSOR_FLAGS(1 << i)
		if e&mask == mask {
			names = append(names, labels_AIRSPEED_SENSOR_FLAGS[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *AIRSPEED_SENSOR_FLAGS) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask AIRSPEED_SENSOR_FLAGS
	for _, label := range labels {
		if value, ok := values_AIRSPEED_SENSOR_FLAGS[label]; ok {
			mask |= value
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e AIRSPEED_SENSOR_FLAGS) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
