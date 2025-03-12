//autogenerated:yes
//nolint:revive,govet,errcheck
package minimal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_TYPE_GENERIC.MarshalText()
		require.NoError(t, err)

		var dec MAV_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_TYPE_GENERIC, dec)
	})
}
