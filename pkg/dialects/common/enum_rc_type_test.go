//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_RC_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e RC_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := RC_TYPE_SPEKTRUM.MarshalText()
		require.NoError(t, err)

		var dec RC_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, RC_TYPE_SPEKTRUM, dec)
	})
}
