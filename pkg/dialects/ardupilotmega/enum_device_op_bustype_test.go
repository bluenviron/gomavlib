//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_DEVICE_OP_BUSTYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e DEVICE_OP_BUSTYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := DEVICE_OP_BUSTYPE_I2C.MarshalText()
		require.NoError(t, err)

		var dec DEVICE_OP_BUSTYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, DEVICE_OP_BUSTYPE_I2C, dec)
	})
}
