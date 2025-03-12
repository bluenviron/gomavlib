//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_CELLULAR_STATUS_FLAG(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e CELLULAR_STATUS_FLAG
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := CELLULAR_STATUS_FLAG_UNKNOWN.MarshalText()
		require.NoError(t, err)

		var dec CELLULAR_STATUS_FLAG
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, CELLULAR_STATUS_FLAG_UNKNOWN, dec)
	})
}
