//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_COLLISION_SRC(t *testing.T) {
	var e MAV_COLLISION_SRC
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}