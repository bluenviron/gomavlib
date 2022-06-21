package dialect

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/message"
)

type (
	MAV_TYPE      uint32 //nolint:revive
	MAV_AUTOPILOT uint32 //nolint:revive
	MAV_MODE_FLAG uint32 //nolint:revive
	MAV_STATE     uint32 //nolint:revive
)

type MessageHeartbeat struct {
	Type           MAV_TYPE      `mavenum:"uint8"`
	Autopilot      MAV_AUTOPILOT `mavenum:"uint8"`
	BaseMode       MAV_MODE_FLAG `mavenum:"uint8"`
	CustomMode     uint32
	SystemStatus   MAV_STATE `mavenum:"uint8"`
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetID() uint32 {
	return 0
}

func TestReadWriter(t *testing.T) {
	_, err := NewReadWriter(&Dialect{3, []message.Message{&MessageHeartbeat{}}})
	require.NoError(t, err)
}
