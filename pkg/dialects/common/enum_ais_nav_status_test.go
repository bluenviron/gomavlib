//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_AIS_NAV_STATUS(t *testing.T) {
	var e AIS_NAV_STATUS
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
