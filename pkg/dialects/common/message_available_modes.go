//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Information about a flight mode.
// The message can be enumerated to get information for all modes, or requested for a particular mode, using MAV_CMD_REQUEST_MESSAGE.
// Specify 0 in param2 to request that the message is emitted for all available modes or the specific index for just one mode.
// The modes must be available/settable for the current vehicle/frame type.
// Each mode should only be emitted once (even if it is both standard and custom).
// Note that the current mode should be emitted in CURRENT_MODE, and that if the mode list can change then AVAILABLE_MODES_MONITOR must be emitted on first change and subsequently streamed.
// See https://mavlink.io/en/services/standard_modes.html
type MessageAvailableModes struct {
	// The total number of available modes for the current vehicle type.
	NumberModes uint8
	// The current mode index within number_modes, indexed from 1. The index is not guaranteed to be persistent, and may change between reboots or if the set of modes change.
	ModeIndex uint8
	// Standard mode.
	StandardMode MAV_STANDARD_MODE `mavenum:"uint8"`
	// A bitfield for use for autopilot-specific flags
	CustomMode uint32
	// Mode properties.
	Properties MAV_MODE_PROPERTY `mavenum:"uint32"`
	// Name of custom mode, with null termination character. Should be omitted for standard modes.
	ModeName string `mavlen:"35"`
}

// GetID implements the message.Message interface.
func (*MessageAvailableModes) GetID() uint32 {
	return 435
}
