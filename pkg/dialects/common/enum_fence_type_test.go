//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_FENCE_TYPE(t *testing.T) {
	var e FENCE_TYPE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
