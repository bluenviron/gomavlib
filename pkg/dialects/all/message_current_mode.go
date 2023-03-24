//autogenerated:yes
//nolint:revive,misspell,govet,lll
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/development"
)

// Get the current mode.
// This should be emitted on any mode change, and broadcast at low rate (nominally 0.5 Hz).
// It may be requested using MAV_CMD_REQUEST_MESSAGE.
type MessageCurrentMode = development.MessageCurrentMode
