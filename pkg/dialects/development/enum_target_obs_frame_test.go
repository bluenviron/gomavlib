//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_TARGET_OBS_FRAME(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e TARGET_OBS_FRAME
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := TARGET_OBS_FRAME_LOCAL_NED.MarshalText()
		require.NoError(t, err)

		var dec TARGET_OBS_FRAME
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, TARGET_OBS_FRAME_LOCAL_NED, dec)
	})
}
