//autogenerated:yes
//nolint:revive,govet,errcheck
package ardupilotmega

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_CAMERA_FEEDBACK_FLAGS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e CAMERA_FEEDBACK_FLAGS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := CAMERA_FEEDBACK_PHOTO.MarshalText()
		require.NoError(t, err)

		var dec CAMERA_FEEDBACK_FLAGS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, CAMERA_FEEDBACK_PHOTO, dec)
	})
}
