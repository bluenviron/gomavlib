//autogenerated:yes
//nolint:revive,govet,errcheck
package csairlink

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_AIRLINK_EYE_HOLE_PUSH_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e AIRLINK_EYE_HOLE_PUSH_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := AIRLINK_HP_NOT_PENETRATED.MarshalText()
		require.NoError(t, err)

		var dec AIRLINK_EYE_HOLE_PUSH_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, AIRLINK_HP_NOT_PENETRATED, dec)
	})
}
