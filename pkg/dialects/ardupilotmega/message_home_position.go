//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Contains the home position.
// The home position is the default position that the system will return to and land on.
// The position must be set automatically by the system during the takeoff, and may also be explicitly set using MAV_CMD_DO_SET_HOME.
// The global and local positions encode the position in the respective coordinate frames, while the q parameter encodes the orientation of the surface.
// Under normal conditions it describes the heading and terrain slope, which can be used by the aircraft to adjust the approach.
// The approach 3D vector describes the point to which the system should fly in normal flight mode and then perform a landing sequence along the vector.
// Note: this message can be requested by sending the MAV_CMD_REQUEST_MESSAGE with param1=242 (or the deprecated MAV_CMD_GET_HOME_POSITION command).
type MessageHomePosition = common.MessageHomePosition
