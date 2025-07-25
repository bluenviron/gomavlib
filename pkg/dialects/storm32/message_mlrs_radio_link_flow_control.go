//autogenerated:yes
//nolint:revive,misspell,govet,lll
package storm32

// Injected by a radio link endpoint into the MAVLink stream for purposes of flow control. Should be emitted only by components with component id MAV_COMP_ID_TELEMETRY_RADIO.
type MessageMlrsRadioLinkFlowControl struct {
	// Transmitted bytes per second, UINT16_MAX: invalid/unknown.
	TxSerRate uint16
	// Received bytes per second, UINT16_MAX: invalid/unknown.
	RxSerRate uint16
	// Transmit bandwidth consumption. Values: 0..100, UINT8_MAX: invalid/unknown.
	TxUsedSerBandwidth uint8
	// Receive bandwidth consumption. Values: 0..100, UINT8_MAX: invalid/unknown.
	RxUsedSerBandwidth uint8
	// For compatibility with legacy method. UINT8_MAX: unknown.
	Txbuf uint8
}

// GetID implements the message.Message interface.
func (*MessageMlrsRadioLinkFlowControl) GetID() uint32 {
	return 60047
}
