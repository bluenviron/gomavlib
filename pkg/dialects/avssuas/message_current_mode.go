//autogenerated:yes
//nolint:revive,misspell,govet,lll
package avssuas

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Get the current mode.
// This should be emitted on any mode change, and broadcast at low rate (nominally 0.5 Hz).
// It may be requested using MAV_CMD_REQUEST_MESSAGE.
// See https://mavlink.io/en/services/standard_modes.html
type MessageCurrentMode = common.MessageCurrentMode
