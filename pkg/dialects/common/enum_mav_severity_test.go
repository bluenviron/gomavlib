//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAV_SEVERITY(t *testing.T) {
	var e MAV_SEVERITY
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
