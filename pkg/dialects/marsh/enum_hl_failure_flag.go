//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Flags to report failure cases over the high latency telemetry.
type HL_FAILURE_FLAG = common.HL_FAILURE_FLAG

const (
	// GPS failure.
	HL_FAILURE_FLAG_GPS HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_GPS
	// Differential pressure sensor failure.
	HL_FAILURE_FLAG_DIFFERENTIAL_PRESSURE HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_DIFFERENTIAL_PRESSURE
	// Absolute pressure sensor failure.
	HL_FAILURE_FLAG_ABSOLUTE_PRESSURE HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_ABSOLUTE_PRESSURE
	// Accelerometer sensor failure.
	HL_FAILURE_FLAG_3D_ACCEL HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_3D_ACCEL
	// Gyroscope sensor failure.
	HL_FAILURE_FLAG_3D_GYRO HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_3D_GYRO
	// Magnetometer sensor failure.
	HL_FAILURE_FLAG_3D_MAG HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_3D_MAG
	// Terrain subsystem failure.
	HL_FAILURE_FLAG_TERRAIN HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_TERRAIN
	// Battery failure/critical low battery.
	HL_FAILURE_FLAG_BATTERY HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_BATTERY
	// RC receiver failure/no RC connection.
	HL_FAILURE_FLAG_RC_RECEIVER HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_RC_RECEIVER
	// Offboard link failure.
	HL_FAILURE_FLAG_OFFBOARD_LINK HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_OFFBOARD_LINK
	// Engine failure.
	HL_FAILURE_FLAG_ENGINE HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_ENGINE
	// Geofence violation.
	HL_FAILURE_FLAG_GEOFENCE HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_GEOFENCE
	// Estimator failure, for example measurement rejection or large variances.
	HL_FAILURE_FLAG_ESTIMATOR HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_ESTIMATOR
	// Mission failure.
	HL_FAILURE_FLAG_MISSION HL_FAILURE_FLAG = common.HL_FAILURE_FLAG_MISSION
)
