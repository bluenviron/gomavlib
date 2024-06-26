//autogenerated:yes
//nolint:revive,misspell,govet,lll
package development

// Information about key components of GNSS receivers, like signal authentication, interference and system errors.
type MessageGnssIntegrity struct {
	// GNSS receiver id. Must match instance ids of other messages from same receiver.
	Id uint8
	// Errors in the GPS system.
	SystemErrors GPS_SYSTEM_ERROR_FLAGS `mavenum:"uint32"`
	// Signal authentication state of the GPS system.
	AuthenticationState GPS_AUTHENTICATION_STATE `mavenum:"uint8"`
	// Signal jamming state of the GPS system.
	JammingState GPS_JAMMING_STATE `mavenum:"uint8"`
	// Signal spoofing state of the GPS system.
	SpoofingState GPS_SPOOFING_STATE `mavenum:"uint8"`
	// The state of the RAIM processing.
	RaimState GPS_RAIM_STATE `mavenum:"uint8"`
	// Horizontal expected accuracy using satellites successfully validated using RAIM.
	RaimHfom uint16
	// Vertical expected accuracy using satellites successfully validated using RAIM.
	RaimVfom uint16
	// An abstract value representing the estimated quality of incoming corrections, or 255 if not available.
	CorrectionsQuality uint8
	// An abstract value representing the overall status of the receiver, or 255 if not available.
	SystemStatusSummary uint8
	// An abstract value representing the quality of incoming GNSS signals, or 255 if not available.
	GnssSignalQuality uint8
	// An abstract value representing the estimated PPK quality, or 255 if not available.
	PostProcessingQuality uint8
}

// GetID implements the message.Message interface.
func (*MessageGnssIntegrity) GetID() uint32 {
	return 441
}
