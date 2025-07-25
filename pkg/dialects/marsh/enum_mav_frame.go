//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Coordinate frames used by MAVLink. Not all frames are supported by all commands, messages, or vehicles.
// Global frames use the following naming conventions:
// - "GLOBAL": Global coordinate frame with WGS84 latitude/longitude and altitude positive over mean sea level (MSL) by default.
// The following modifiers may be used with "GLOBAL":
// - "RELATIVE_ALT": Altitude is relative to the vehicle home position rather than MSL.
// - "TERRAIN_ALT": Altitude is relative to ground level rather than MSL.
// - "INT": Latitude/longitude (in degrees) are scaled by multiplying by 1E7.
// Local frames use the following naming conventions:
// - "LOCAL": Origin of local frame is fixed relative to earth. Unless otherwise specified this origin is the origin of the vehicle position-estimator ("EKF").
// - "BODY": Origin of local frame travels with the vehicle. NOTE, "BODY" does NOT indicate alignment of frame axis with vehicle attitude.
// - "OFFSET": Deprecated synonym for "BODY" (origin travels with the vehicle). Not to be used for new frames.
// Some deprecated frames do not follow these conventions (e.g. MAV_FRAME_BODY_NED and MAV_FRAME_BODY_OFFSET_NED).
type MAV_FRAME = common.MAV_FRAME

const (
	// Global (WGS84) coordinate frame + altitude relative to mean sea level (MSL).
	MAV_FRAME_GLOBAL MAV_FRAME = common.MAV_FRAME_GLOBAL
	// NED local tangent frame (x: North, y: East, z: Down) with origin fixed relative to earth.
	MAV_FRAME_LOCAL_NED MAV_FRAME = common.MAV_FRAME_LOCAL_NED
	// NOT a coordinate frame, indicates a mission command.
	MAV_FRAME_MISSION MAV_FRAME = common.MAV_FRAME_MISSION
	// Global (WGS84) coordinate frame + altitude relative to the home position.
	MAV_FRAME_GLOBAL_RELATIVE_ALT MAV_FRAME = common.MAV_FRAME_GLOBAL_RELATIVE_ALT
	// ENU local tangent frame (x: East, y: North, z: Up) with origin fixed relative to earth.
	MAV_FRAME_LOCAL_ENU MAV_FRAME = common.MAV_FRAME_LOCAL_ENU
	// Global (WGS84) coordinate frame (scaled) + altitude relative to mean sea level (MSL).
	MAV_FRAME_GLOBAL_INT MAV_FRAME = common.MAV_FRAME_GLOBAL_INT
	// Global (WGS84) coordinate frame (scaled) + altitude relative to the home position.
	MAV_FRAME_GLOBAL_RELATIVE_ALT_INT MAV_FRAME = common.MAV_FRAME_GLOBAL_RELATIVE_ALT_INT
	// NED local tangent frame (x: North, y: East, z: Down) with origin that travels with the vehicle.
	MAV_FRAME_LOCAL_OFFSET_NED MAV_FRAME = common.MAV_FRAME_LOCAL_OFFSET_NED
	// Same as MAV_FRAME_LOCAL_NED when used to represent position values. Same as MAV_FRAME_BODY_FRD when used with velocity/acceleration values.
	MAV_FRAME_BODY_NED MAV_FRAME = common.MAV_FRAME_BODY_NED
	// This is the same as MAV_FRAME_BODY_FRD.
	MAV_FRAME_BODY_OFFSET_NED MAV_FRAME = common.MAV_FRAME_BODY_OFFSET_NED
	// Global (WGS84) coordinate frame with AGL altitude (altitude at ground level).
	MAV_FRAME_GLOBAL_TERRAIN_ALT MAV_FRAME = common.MAV_FRAME_GLOBAL_TERRAIN_ALT
	// Global (WGS84) coordinate frame (scaled) with AGL altitude (altitude at ground level).
	MAV_FRAME_GLOBAL_TERRAIN_ALT_INT MAV_FRAME = common.MAV_FRAME_GLOBAL_TERRAIN_ALT_INT
	// FRD local frame aligned to the vehicle's attitude (x: Forward, y: Right, z: Down) with an origin that travels with vehicle.
	MAV_FRAME_BODY_FRD MAV_FRAME = common.MAV_FRAME_BODY_FRD
	// MAV_FRAME_BODY_FLU - Body fixed frame of reference, Z-up (x: Forward, y: Left, z: Up).
	MAV_FRAME_RESERVED_13 MAV_FRAME = common.MAV_FRAME_RESERVED_13
	// MAV_FRAME_MOCAP_NED - Odometry local coordinate frame of data given by a motion capture system, Z-down (x: North, y: East, z: Down).
	MAV_FRAME_RESERVED_14 MAV_FRAME = common.MAV_FRAME_RESERVED_14
	// MAV_FRAME_MOCAP_ENU - Odometry local coordinate frame of data given by a motion capture system, Z-up (x: East, y: North, z: Up).
	MAV_FRAME_RESERVED_15 MAV_FRAME = common.MAV_FRAME_RESERVED_15
	// MAV_FRAME_VISION_NED - Odometry local coordinate frame of data given by a vision estimation system, Z-down (x: North, y: East, z: Down).
	MAV_FRAME_RESERVED_16 MAV_FRAME = common.MAV_FRAME_RESERVED_16
	// MAV_FRAME_VISION_ENU - Odometry local coordinate frame of data given by a vision estimation system, Z-up (x: East, y: North, z: Up).
	MAV_FRAME_RESERVED_17 MAV_FRAME = common.MAV_FRAME_RESERVED_17
	// MAV_FRAME_ESTIM_NED - Odometry local coordinate frame of data given by an estimator running onboard the vehicle, Z-down (x: North, y: East, z: Down).
	MAV_FRAME_RESERVED_18 MAV_FRAME = common.MAV_FRAME_RESERVED_18
	// MAV_FRAME_ESTIM_ENU - Odometry local coordinate frame of data given by an estimator running onboard the vehicle, Z-up (x: East, y: North, z: Up).
	MAV_FRAME_RESERVED_19 MAV_FRAME = common.MAV_FRAME_RESERVED_19
	// FRD local tangent frame (x: Forward, y: Right, z: Down) with origin fixed relative to earth. The forward axis is aligned to the front of the vehicle in the horizontal plane.
	MAV_FRAME_LOCAL_FRD MAV_FRAME = common.MAV_FRAME_LOCAL_FRD
	// FLU local tangent frame (x: Forward, y: Left, z: Up) with origin fixed relative to earth. The forward axis is aligned to the front of the vehicle in the horizontal plane.
	MAV_FRAME_LOCAL_FLU MAV_FRAME = common.MAV_FRAME_LOCAL_FLU
)
