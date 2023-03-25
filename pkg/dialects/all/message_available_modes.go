//autogenerated:yes
//nolint:revive,misspell,govet,lll
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/development"
)

// Get information about a particular flight modes.
// The message can be enumerated or requested for a particular mode using MAV_CMD_REQUEST_MESSAGE.
// Specify 0 in param2 to request that the message is emitted for all available modes or the specific index for just one mode.
// The modes must be available/settable for the current vehicle/frame type.
// Each modes should only be emitted once (even if it is both standard and custom).
type MessageAvailableModes = development.MessageAvailableModes
