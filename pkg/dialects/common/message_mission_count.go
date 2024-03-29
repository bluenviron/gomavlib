//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// This message is emitted as response to MISSION_REQUEST_LIST by the MAV and to initiate a write transaction. The GCS can then request the individual mission item based on the knowledge of the total number of waypoints.
type MessageMissionCount struct {
	// System ID
	TargetSystem uint8
	// Component ID
	TargetComponent uint8
	// Number of mission items in the sequence
	Count uint16
	// Mission type.
	MissionType MAV_MISSION_TYPE `mavenum:"uint8" mavext:"true"`
	// Id of current on-vehicle mission, fence, or rally point plan (on download from vehicle).
	// This field is used when downloading a plan from a vehicle to a GCS.
	// 0 on upload to the vehicle from GCS.
	// 0 if plan ids are not supported.
	// The current on-vehicle plan ids are streamed in `MISSION_CURRENT`, allowing a GCS to determine if any part of the plan has changed and needs to be re-uploaded.
	// The ids are recalculated by the vehicle when any part of the on-vehicle plan changes (when a new plan is uploaded, the vehicle returns the new id to the GCS in MISSION_ACK).
	OpaqueId uint32 `mavext:"true"`
}

// GetID implements the message.Message interface.
func (*MessageMissionCount) GetID() uint32 {
	return 44
}
