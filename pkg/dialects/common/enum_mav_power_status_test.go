//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_POWER_STATUS(t *testing.T) {
	var e MAV_POWER_STATUS
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
