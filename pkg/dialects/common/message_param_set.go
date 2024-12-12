//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Set a parameter value (write new value to permanent storage).
// The receiving component should acknowledge the new parameter value by broadcasting a PARAM_VALUE message (broadcasting ensures that multiple GCS all have an up-to-date list of all parameters). If the sending GCS did not receive a PARAM_VALUE within its timeout time, it should re-send the PARAM_SET message. The parameter microservice is documented at https://mavlink.io/en/services/parameter.html.
type MessageParamSet struct {
	// System ID
	TargetSystem uint8
	// Component ID
	TargetComponent uint8
	// Onboard parameter id, terminated by NULL if the length is less than 16 human-readable chars and WITHOUT null termination (NULL) byte if the length is exactly 16 chars - applications have to provide 16+1 bytes storage if the ID is stored as string
	ParamId string `mavlen:"16"`
	// Onboard parameter value
	ParamValue float32
	// Onboard parameter type.
	ParamType MAV_PARAM_TYPE `mavenum:"uint8"`
}

// GetID implements the message.Message interface.
func (*MessageParamSet) GetID() uint32 {
	return 23
}
