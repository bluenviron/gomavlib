//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_GOPRO_REQUEST_STATUS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e GOPRO_REQUEST_STATUS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := GOPRO_REQUEST_SUCCESS.MarshalText()
		require.NoError(t, err)

		var dec GOPRO_REQUEST_STATUS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, GOPRO_REQUEST_SUCCESS, dec)
	})
}
