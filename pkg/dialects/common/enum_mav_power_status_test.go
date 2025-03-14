//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_POWER_STATUS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_POWER_STATUS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_POWER_STATUS_BRICK_VALID.MarshalText()
		require.NoError(t, err)

		var dec MAV_POWER_STATUS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_POWER_STATUS_BRICK_VALID, dec)
	})
}
