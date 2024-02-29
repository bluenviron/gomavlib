package dialect

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

type (
	MAV_TYPE      uint64 //nolint:revive
	MAV_AUTOPILOT uint64 //nolint:revive
	MAV_MODE_FLAG uint64 //nolint:revive
	MAV_STATE     uint64 //nolint:revive
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

type Invalid struct{}

func (*Invalid) GetID() uint32 {
	return 0
}

func TestReadWriter(t *testing.T) {
	rw, err := NewReadWriter(&Dialect{3, []message.Message{&MessageHeartbeat{}}})
	require.NoError(t, err)

	mrw := rw.GetMessage(0)
	require.NotNil(t, mrw)

	mrw = rw.GetMessage(1)
	require.Nil(t, mrw)
}

func TestReadWriterErrors(t *testing.T) {
	for _, ca := range []struct {
		name    string
		dialect *Dialect
		err     string
	}{
		{
			"duplicate message",
			&Dialect{3, []message.Message{
				&MessageHeartbeat{},
				&MessageHeartbeat{},
			}},
			"duplicate message with id 0",
		},
		{
			"invalid message",
			&Dialect{3, []message.Message{
				&Invalid{},
			}},
			"message *dialect.Invalid: struct name must begin with 'Message'",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			_, err := NewReadWriter(ca.dialect)
			require.EqualError(t, err, ca.err)
		})
	}
}
