//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_COLLISION_SRC(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_COLLISION_SRC
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_COLLISION_SRC_ADSB.MarshalText()
		require.NoError(t, err)

		var dec MAV_COLLISION_SRC
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_COLLISION_SRC_ADSB, dec)
	})
}
