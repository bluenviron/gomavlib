//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_ACTUATOR_OUTPUT_FUNCTION(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e ACTUATOR_OUTPUT_FUNCTION
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := ACTUATOR_OUTPUT_FUNCTION_NONE.MarshalText()
		require.NoError(t, err)

		var dec ACTUATOR_OUTPUT_FUNCTION
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, ACTUATOR_OUTPUT_FUNCTION_NONE, dec)
	})
}
