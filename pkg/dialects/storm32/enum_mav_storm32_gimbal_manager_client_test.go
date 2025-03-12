//autogenerated:yes
//nolint:revive,govet,errcheck
package storm32

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_STORM32_GIMBAL_MANAGER_CLIENT(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_STORM32_GIMBAL_MANAGER_CLIENT
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_STORM32_GIMBAL_MANAGER_CLIENT_NONE.MarshalText()
		require.NoError(t, err)

		var dec MAV_STORM32_GIMBAL_MANAGER_CLIENT
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_STORM32_GIMBAL_MANAGER_CLIENT_NONE, dec)
	})
}
