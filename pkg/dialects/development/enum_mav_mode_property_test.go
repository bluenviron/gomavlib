//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"
)

func TestEnum_MAV_MODE_PROPERTY(t *testing.T) {
	var e MAV_MODE_PROPERTY
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}