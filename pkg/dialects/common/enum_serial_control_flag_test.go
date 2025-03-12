//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_SERIAL_CONTROL_FLAG(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e SERIAL_CONTROL_FLAG
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := SERIAL_CONTROL_FLAG_REPLY.MarshalText()
		require.NoError(t, err)

		var dec SERIAL_CONTROL_FLAG
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, SERIAL_CONTROL_FLAG_REPLY, dec)
	})
}
