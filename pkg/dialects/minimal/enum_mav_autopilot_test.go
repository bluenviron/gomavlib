//autogenerated:yes
//nolint:revive,govet,errcheck
package minimal

import (
	"testing"
)

func TestEnum_MAV_AUTOPILOT(t *testing.T) {
	var e MAV_AUTOPILOT
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
