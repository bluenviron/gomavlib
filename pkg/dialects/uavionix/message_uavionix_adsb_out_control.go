//autogenerated:yes
//nolint:revive,misspell,govet,lll
package uavionix

// Control message with all data sent in UCP control message.
type MessageUavionixAdsbOutControl struct {
	// ADS-B transponder control state flags
	State UAVIONIX_ADSB_OUT_CONTROL_STATE `mavenum:"uint8"`
	// Barometric pressure altitude (MSL) relative to a standard atmosphere of 1013.2 mBar and NOT bar corrected altitude (m * 1E-3). (up +ve). If unknown set to INT32_MAX
	Baroaltmsl int32 `mavname:"baroAltMSL"`
	// Mode A code (typically 1200 [0x04B0] for VFR)
	Squawk uint16
	// Emergency status
	Emergencystatus UAVIONIX_ADSB_EMERGENCY_STATUS `mavenum:"uint8" mavname:"emergencyStatus"`
	// Flight Identification: 8 ASCII characters, '0' through '9', 'A' through 'Z' or space. Spaces (0x20) used as a trailing pad character, or when call sign is unavailable.
	FlightId string `mavlen:"8"`
	// X-Bit enable (military transponders only)
	XBit UAVIONIX_ADSB_XBIT `mavenum:"uint8"`
}

// GetID implements the message.Message interface.
func (*MessageUavionixAdsbOutControl) GetID() uint32 {
	return 10007
}
