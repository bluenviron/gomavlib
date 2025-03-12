//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_PREFLIGHT_STORAGE_MISSION_ACTION(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e PREFLIGHT_STORAGE_MISSION_ACTION
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MISSION_READ_PERSISTENT.MarshalText()
		require.NoError(t, err)

		var dec PREFLIGHT_STORAGE_MISSION_ACTION
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MISSION_READ_PERSISTENT, dec)
	})
}
