//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_ESTIMATOR_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_ESTIMATOR_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_ESTIMATOR_TYPE_UNKNOWN.MarshalText()
		require.NoError(t, err)

		var dec MAV_ESTIMATOR_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_ESTIMATOR_TYPE_UNKNOWN, dec)
	})
}
