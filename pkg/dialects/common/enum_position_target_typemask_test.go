//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_POSITION_TARGET_TYPEMASK(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e POSITION_TARGET_TYPEMASK
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := POSITION_TARGET_TYPEMASK_X_IGNORE.MarshalText()
		require.NoError(t, err)

		var dec POSITION_TARGET_TYPEMASK
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, POSITION_TARGET_TYPEMASK_X_IGNORE, dec)
	})
}
