//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_CAMERA_TRACKING_MODE(t *testing.T) {
	var e CAMERA_TRACKING_MODE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
