//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_AUTOTUNE_AXIS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e AUTOTUNE_AXIS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := AUTOTUNE_AXIS_DEFAULT.MarshalText()
		require.NoError(t, err)

		var dec AUTOTUNE_AXIS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, AUTOTUNE_AXIS_DEFAULT, dec)
	})
}
