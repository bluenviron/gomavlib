//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_STORAGE_STATUS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e STORAGE_STATUS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := STORAGE_STATUS_EMPTY.MarshalText()
		require.NoError(t, err)

		var dec STORAGE_STATUS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, STORAGE_STATUS_EMPTY, dec)
	})
}
