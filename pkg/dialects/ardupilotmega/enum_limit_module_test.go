//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_LIMIT_MODULE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e LIMIT_MODULE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := LIMIT_GPSLOCK.MarshalText()
		require.NoError(t, err)

		var dec LIMIT_MODULE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, LIMIT_GPSLOCK, dec)
	})
}
