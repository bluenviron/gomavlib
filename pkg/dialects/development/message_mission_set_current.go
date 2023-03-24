//autogenerated:yes
//nolint:revive,misspell,govet,lll
package development

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Set the mission item with sequence number seq as the current item and emit MISSION_CURRENT (whether or not the mission number changed).
// If a mission is currently being executed, the system will continue to this new mission item on the shortest path, skipping any intermediate mission items.
// Note that mission jump repeat counters are not reset (see MAV_CMD_DO_JUMP param2).
// This message may trigger a mission state-machine change on some systems: for example from MISSION_STATE_NOT_STARTED or MISSION_STATE_PAUSED to MISSION_STATE_ACTIVE.
// If the system is in mission mode, on those systems this command might therefore start, restart or resume the mission.
// If the system is not in mission mode this message must not trigger a switch to mission mode.
type MessageMissionSetCurrent = common.MessageMissionSetCurrent
