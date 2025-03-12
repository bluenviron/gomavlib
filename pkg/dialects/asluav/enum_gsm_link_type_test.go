//autogenerated:yes
//nolint:revive,govet,errcheck
package asluav

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_GSM_LINK_TYPE(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e GSM_LINK_TYPE
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := GSM_LINK_TYPE_NONE.MarshalText()
		require.NoError(t, err)

		var dec GSM_LINK_TYPE
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, GSM_LINK_TYPE_NONE, dec)
	})
}
