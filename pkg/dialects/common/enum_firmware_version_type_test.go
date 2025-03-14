//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_FIRMWARE_VERSION_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e FIRMWARE_VERSION_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := FIRMWARE_VERSION_TYPE_DEV.MarshalText()
		require.NoError(t, err)

		var dec FIRMWARE_VERSION_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, FIRMWARE_VERSION_TYPE_DEV, dec)
	})
}
