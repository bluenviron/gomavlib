//autogenerated:yes
//nolint:revive,govet,errcheck
package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_GRIPPER_ACTIONS(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e GRIPPER_ACTIONS
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := GRIPPER_ACTION_RELEASE.MarshalText()
		require.NoError(t, err)

		var dec GRIPPER_ACTIONS
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, GRIPPER_ACTION_RELEASE, dec)
	})
}
