//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_CELLULAR_NETWORK_RADIO_TYPE(t *testing.T) {
	var e CELLULAR_NETWORK_RADIO_TYPE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}