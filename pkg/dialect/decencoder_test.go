package dialect

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/msg"
)

type (
	MAV_TYPE      int //nolint:revive
	MAV_AUTOPILOT int //nolint:revive
	MAV_MODE_FLAG int //nolint:revive
	MAV_STATE     int //nolint:revive
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

func TestDecEncoder(t *testing.T) {
	_, err := NewDecEncoder(&Dialect{3, []msg.Message{&MessageHeartbeat{}}})
	require.NoError(t, err)
}
