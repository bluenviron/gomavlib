//autogenerated:yes
//nolint:revive,govet,errcheck
package development

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_MAV_CMD(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e MAV_CMD
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := MAV_CMD_NAV_WAYPOINT.MarshalText()
		require.NoError(t, err)

		var dec MAV_CMD
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, MAV_CMD_NAV_WAYPOINT, dec)
	})
}
