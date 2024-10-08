//autogenerated:yes
//nolint:revive,misspell,govet,lll
package common

// Settings of a camera. Can be requested with a MAV_CMD_REQUEST_MESSAGE command.
type MessageCameraSettings struct {
	// Timestamp (time since system boot).
	TimeBootMs uint32
	// Camera mode
	ModeId CAMERA_MODE `mavenum:"uint8"`
	// Current zoom level as a percentage of the full range (0.0 to 100.0, NaN if not known)
	Zoomlevel float32 `mavext:"true" mavname:"zoomLevel"`
	// Current focus level as a percentage of the full range (0.0 to 100.0, NaN if not known)
	Focuslevel float32 `mavext:"true" mavname:"focusLevel"`
	// Camera id of a non-MAVLink camera attached to an autopilot (1-6).  0 if the component is a MAVLink camera (with its own component id).
	CameraDeviceId uint8 `mavext:"true"`
}

// GetID implements the message.Message interface.
func (*MessageCameraSettings) GetID() uint32 {
	return 260
}
