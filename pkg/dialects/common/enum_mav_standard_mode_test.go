//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_STANDARD_MODE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_STANDARD_MODE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_STANDARD_MODE_NON_STANDARD.MarshalText()
		require.NoError(t, err)

		var dec MAV_STANDARD_MODE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_STANDARD_MODE_NON_STANDARD, dec)
	})
}
