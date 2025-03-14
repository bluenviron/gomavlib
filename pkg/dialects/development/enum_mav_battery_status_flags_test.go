//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_BATTERY_STATUS_FLAGS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_BATTERY_STATUS_FLAGS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_BATTERY_STATUS_FLAGS_NOT_READY_TO_USE.MarshalText()
		require.NoError(t, err)

		var dec MAV_BATTERY_STATUS_FLAGS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_BATTERY_STATUS_FLAGS_NOT_READY_TO_USE, dec)
	})
}
