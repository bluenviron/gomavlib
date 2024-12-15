//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Possible actions an aircraft can take to avoid a collision.
type MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION

const (
	// Ignore any potential collisions
	MAV_COLLISION_ACTION_NONE MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_NONE
	// Report potential collision
	MAV_COLLISION_ACTION_REPORT MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_REPORT
	// Ascend or Descend to avoid threat
	MAV_COLLISION_ACTION_ASCEND_OR_DESCEND MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_ASCEND_OR_DESCEND
	// Move horizontally to avoid threat
	MAV_COLLISION_ACTION_MOVE_HORIZONTALLY MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_MOVE_HORIZONTALLY
	// Aircraft to move perpendicular to the collision's velocity vector
	MAV_COLLISION_ACTION_MOVE_PERPENDICULAR MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_MOVE_PERPENDICULAR
	// Aircraft to fly directly back to its launch point
	MAV_COLLISION_ACTION_RTL MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_RTL
	// Aircraft to stop in place
	MAV_COLLISION_ACTION_HOVER MAV_COLLISION_ACTION = common.MAV_COLLISION_ACTION_HOVER
)
