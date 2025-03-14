//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_CELLULAR_NETWORK_FAILED_REASON(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e CELLULAR_NETWORK_FAILED_REASON
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := CELLULAR_NETWORK_FAILED_REASON_NONE.MarshalText()
		require.NoError(t, err)

		var dec CELLULAR_NETWORK_FAILED_REASON
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, CELLULAR_NETWORK_FAILED_REASON_NONE, dec)
	})
}
