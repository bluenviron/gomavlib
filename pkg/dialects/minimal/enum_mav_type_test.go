//autogenerated:yes
//nolint:revive,govet,errcheck
package minimal

import (
	"testing"
)

func TestEnum_MAV_TYPE(t *testing.T) {
	var e MAV_TYPE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
