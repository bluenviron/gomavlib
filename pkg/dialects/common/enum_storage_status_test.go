//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_STORAGE_STATUS(t *testing.T) {
	var e STORAGE_STATUS
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}