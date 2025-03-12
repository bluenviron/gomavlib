//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_LED_CONTROL_PATTERN(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e LED_CONTROL_PATTERN
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := LED_CONTROL_PATTERN_OFF.MarshalText()
		require.NoError(t, err)

		var dec LED_CONTROL_PATTERN
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, LED_CONTROL_PATTERN_OFF, dec)
	})
}
