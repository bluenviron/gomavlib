//autogenerated:yes
//nolint:revive,misspell,govet,lll
package marsh

// Data for tracking of pilot eye gaze. This is the primary message for MARSH_TYPE_EYE_TRACKER.
type MessageEyeTrackingData struct {
	// Timestamp (time since system boot).
	TimeUsec uint64
	// Sensor ID, used for identifying the device and/or person tracked. Set to zero if unknown/unused.
	SensorId uint8
	// X axis of gaze origin point, NaN if unknown. The reference system depends on specific application.
	GazeOriginX float32
	// Y axis of gaze origin point, NaN if unknown. The reference system depends on specific application.
	GazeOriginY float32
	// Z axis of gaze origin point, NaN if unknown. The reference system depends on specific application.
	GazeOriginZ float32
	// X axis of gaze direction vector, expected to be normalized to unit magnitude, NaN if unknown. The reference system should match origin point.
	GazeDirectionX float32
	// Y axis of gaze direction vector, expected to be normalized to unit magnitude, NaN if unknown. The reference system should match origin point.
	GazeDirectionY float32
	// Z axis of gaze direction vector, expected to be normalized to unit magnitude, NaN if unknown. The reference system should match origin point.
	GazeDirectionZ float32
	// Gaze focal point on video feed x value (normalized 0..1, 0 is left, 1 is right), NaN if unknown
	VideoGazeX float32
	// Gaze focal point on video feed y value (normalized 0..1, 0 is top, 1 is bottom), NaN if unknown
	VideoGazeY float32
	// Identifier of surface for 2D gaze point, or an identified region when surface point is invalid. Set to zero if unknown/unused.
	SurfaceId uint8
	// Gaze focal point on surface x value (normalized 0..1, 0 is left, 1 is right), NaN if unknown
	SurfaceGazeX float32
	// Gaze focal point on surface y value (normalized 0..1, 0 is top, 1 is bottom), NaN if unknown
	SurfaceGazeY float32
}

// GetID implements the message.Message interface.
func (*MessageEyeTrackingData) GetID() uint32 {
	return 52505
}
