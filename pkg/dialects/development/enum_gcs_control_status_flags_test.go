//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_GCS_CONTROL_STATUS_FLAGS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e GCS_CONTROL_STATUS_FLAGS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := GCS_CONTROL_STATUS_FLAGS_SYSTEM_MANAGER.MarshalText()
		require.NoError(t, err)

		var dec GCS_CONTROL_STATUS_FLAGS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, GCS_CONTROL_STATUS_FLAGS_SYSTEM_MANAGER, dec)
	})
}
