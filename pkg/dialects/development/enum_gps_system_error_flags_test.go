//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_GPS_SYSTEM_ERROR_FLAGS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e GPS_SYSTEM_ERROR_FLAGS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := GPS_SYSTEM_ERROR_INCOMING_CORRECTIONS.MarshalText()
		require.NoError(t, err)

		var dec GPS_SYSTEM_ERROR_FLAGS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, GPS_SYSTEM_ERROR_INCOMING_CORRECTIONS, dec)
	})
}
