//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_LANDED_STATE(t *testing.T) {
	var e MAV_LANDED_STATE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}