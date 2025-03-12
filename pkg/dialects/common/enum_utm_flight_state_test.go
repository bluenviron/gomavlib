//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_UTM_FLIGHT_STATE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e UTM_FLIGHT_STATE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := UTM_FLIGHT_STATE_UNKNOWN.MarshalText()
		require.NoError(t, err)

		var dec UTM_FLIGHT_STATE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, UTM_FLIGHT_STATE_UNKNOWN, dec)
	})
}
