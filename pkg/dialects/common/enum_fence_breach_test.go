//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_FENCE_BREACH(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e FENCE_BREACH
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := FENCE_BREACH_NONE.MarshalText()
		require.NoError(t, err)

		var dec FENCE_BREACH
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, FENCE_BREACH_NONE, dec)
	})
}
