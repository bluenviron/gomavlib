//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_PLANE_MODE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e PLANE_MODE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := PLANE_MODE_MANUAL.MarshalText()
		require.NoError(t, err)

		var dec PLANE_MODE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, PLANE_MODE_MANUAL, dec)
	})
}
