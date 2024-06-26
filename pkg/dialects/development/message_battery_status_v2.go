//autogenerated:yes
//nolint:revive,misspell,govet,lll
package development

// Battery dynamic information.
// This should be streamed (nominally at 1Hz).
// Static/invariant battery information is sent in BATTERY_INFO.
// Note that smart batteries should set the MAV_BATTERY_STATUS_FLAGS_CAPACITY_RELATIVE_TO_FULL bit to indicate that supplied capacity values are relative to a battery that is known to be full.
// Power monitors would not set this bit, indicating that capacity_consumed is relative to drone power-on, and that other values are estimated based on the assumption that the battery was full on power-on.
type MessageBatteryStatusV2 struct {
	// Battery ID
	Id uint8
	// Temperature of the whole battery pack (not internal electronics). INT16_MAX field not provided.
	Temperature int16
	// Battery voltage (total). NaN: field not provided.
	Voltage float32
	// Battery current (through all cells/loads). Positive value when discharging and negative if charging. NaN: field not provided.
	Current float32
	// Consumed charge. NaN: field not provided. This is either the consumption since power-on or since the battery was full, depending on the value of MAV_BATTERY_STATUS_FLAGS_CAPACITY_RELATIVE_TO_FULL.
	CapacityConsumed float32
	// Remaining charge (until empty). NaN: field not provided. Note: If MAV_BATTERY_STATUS_FLAGS_CAPACITY_RELATIVE_TO_FULL is unset, this value is based on the assumption the battery was full when the system was powered.
	CapacityRemaining float32
	// Remaining battery energy. Values: [0-100], UINT8_MAX: field not provided.
	PercentRemaining uint8
	// Fault, health, readiness, and other status indications.
	StatusFlags MAV_BATTERY_STATUS_FLAGS `mavenum:"uint32"`
}

// GetID implements the message.Message interface.
func (*MessageBatteryStatusV2) GetID() uint32 {
	return 369
}
