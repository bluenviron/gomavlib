//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"
)

func TestEnum_MAVLINK_DATA_STREAM_TYPE(t *testing.T) {
	var e MAVLINK_DATA_STREAM_TYPE
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}
