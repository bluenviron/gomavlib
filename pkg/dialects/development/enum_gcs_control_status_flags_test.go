//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"
)

func TestEnum_GCS_CONTROL_STATUS_FLAGS(t *testing.T) {
	var e GCS_CONTROL_STATUS_FLAGS
	e.UnmarshalText([]byte{})
	e.MarshalText()
	e.String()
}