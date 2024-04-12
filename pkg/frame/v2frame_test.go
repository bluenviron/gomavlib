package frame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestV2Frame(t *testing.T) {
	f := &V2Frame{
		SequenceNumber: 123,
		SystemID:       56,
		ComponentID:    89,
		Message:        nil,
		Checksum:       31415,
	}
	require.Equal(t, uint8(56), f.GetSystemID())
	require.Equal(t, uint8(89), f.GetComponentID())
	require.Equal(t, uint8(123), f.GetSequenceNumber())
	require.Equal(t, uint16(31415), f.GetChecksum())
}
