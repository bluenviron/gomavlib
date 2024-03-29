//autogenerated:yes
//nolint:revive,misspell,govet,lll
package development

// Current limits for horizontal speed, vertical speed and yaw rate, as set by SET_VELOCITY_LIMITS.
type MessageVelocityLimits struct {
	// Limit for horizontal movement in MAV_FRAME_LOCAL_NED. NaN: No limit applied
	HorizontalSpeedLimit float32
	// Limit for vertical movement in MAV_FRAME_LOCAL_NED. NaN: No limit applied
	VerticalSpeedLimit float32
	// Limit for vehicle turn rate around its yaw axis. NaN: No limit applied
	YawRateLimit float32
}

// GetID implements the message.Message interface.
func (*MessageVelocityLimits) GetID() uint32 {
	return 355
}
