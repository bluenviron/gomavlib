//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_FAILURE_UNIT(t *testing.T) {
	var e FAILURE_UNIT
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
