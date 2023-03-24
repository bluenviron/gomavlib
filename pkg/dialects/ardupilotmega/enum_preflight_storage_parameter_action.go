//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Actions for reading/writing parameters between persistent and volatile storage when using MAV_CMD_PREFLIGHT_STORAGE.
// (Commonly parameters are loaded from persistent storage (flash/EEPROM) into volatile storage (RAM) on startup and written back when they are changed.)
type PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PREFLIGHT_STORAGE_PARAMETER_ACTION

const (
	// Read all parameters from persistent storage. Replaces values in volatile storage.
	PARAM_READ_PERSISTENT PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PARAM_READ_PERSISTENT
	// Write all parameter values to persistent storage (flash/EEPROM)
	PARAM_WRITE_PERSISTENT PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PARAM_WRITE_PERSISTENT
	// Reset all user configurable parameters to their default value (including airframe selection, sensor calibration data, safety settings, and so on). Does not reset values that contain operation counters and vehicle computed statistics.
	PARAM_RESET_CONFIG_DEFAULT PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PARAM_RESET_CONFIG_DEFAULT
	// Reset only sensor calibration parameters to factory defaults (or firmware default if not available)
	PARAM_RESET_SENSOR_DEFAULT PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PARAM_RESET_SENSOR_DEFAULT
	// Reset all parameters, including operation counters, to default values
	PARAM_RESET_ALL_DEFAULT PREFLIGHT_STORAGE_PARAMETER_ACTION = common.PARAM_RESET_ALL_DEFAULT
)
