//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_ADSB_ALTITUDE_TYPE(t *testing.T) {
	var e ADSB_ALTITUDE_TYPE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
