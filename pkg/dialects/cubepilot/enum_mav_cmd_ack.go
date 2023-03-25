//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package cubepilot

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// ACK / NACK / ERROR values as a result of MAV_CMDs and for mission item transmission.
type MAV_CMD_ACK = common.MAV_CMD_ACK

const (
	// Command / mission item is ok.
	MAV_CMD_ACK_OK MAV_CMD_ACK = common.MAV_CMD_ACK_OK
	// Generic error message if none of the other reasons fails or if no detailed error reporting is implemented.
	MAV_CMD_ACK_ERR_FAIL MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_FAIL
	// The system is refusing to accept this command from this source / communication partner.
	MAV_CMD_ACK_ERR_ACCESS_DENIED MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_ACCESS_DENIED
	// Command or mission item is not supported, other commands would be accepted.
	MAV_CMD_ACK_ERR_NOT_SUPPORTED MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_NOT_SUPPORTED
	// The coordinate frame of this command / mission item is not supported.
	MAV_CMD_ACK_ERR_COORDINATE_FRAME_NOT_SUPPORTED MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_COORDINATE_FRAME_NOT_SUPPORTED
	// The coordinate frame of this command is ok, but he coordinate values exceed the safety limits of this system. This is a generic error, please use the more specific error messages below if possible.
	MAV_CMD_ACK_ERR_COORDINATES_OUT_OF_RANGE MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_COORDINATES_OUT_OF_RANGE
	// The X or latitude value is out of range.
	MAV_CMD_ACK_ERR_X_LAT_OUT_OF_RANGE MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_X_LAT_OUT_OF_RANGE
	// The Y or longitude value is out of range.
	MAV_CMD_ACK_ERR_Y_LON_OUT_OF_RANGE MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_Y_LON_OUT_OF_RANGE
	// The Z or altitude value is out of range.
	MAV_CMD_ACK_ERR_Z_ALT_OUT_OF_RANGE MAV_CMD_ACK = common.MAV_CMD_ACK_ERR_Z_ALT_OUT_OF_RANGE
)
