//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Actions for reading and writing plan information (mission, rally points, geofence) between persistent and volatile storage when using MAV_CMD_PREFLIGHT_STORAGE.
// (Commonly missions are loaded from persistent storage (flash/EEPROM) into volatile storage (RAM) on startup and written back when they are changed.)
type PREFLIGHT_STORAGE_MISSION_ACTION uint32

const (
	// Read current mission data from persistent storage
	MISSION_READ_PERSISTENT PREFLIGHT_STORAGE_MISSION_ACTION = 0
	// Write current mission data to persistent storage
	MISSION_WRITE_PERSISTENT PREFLIGHT_STORAGE_MISSION_ACTION = 1
	// Erase all mission data stored on the vehicle (both persistent and volatile storage)
	MISSION_RESET_DEFAULT PREFLIGHT_STORAGE_MISSION_ACTION = 2
)

var labels_PREFLIGHT_STORAGE_MISSION_ACTION = map[PREFLIGHT_STORAGE_MISSION_ACTION]string{
	MISSION_READ_PERSISTENT:  "MISSION_READ_PERSISTENT",
	MISSION_WRITE_PERSISTENT: "MISSION_WRITE_PERSISTENT",
	MISSION_RESET_DEFAULT:    "MISSION_RESET_DEFAULT",
}

var values_PREFLIGHT_STORAGE_MISSION_ACTION = map[string]PREFLIGHT_STORAGE_MISSION_ACTION{
	"MISSION_READ_PERSISTENT":  MISSION_READ_PERSISTENT,
	"MISSION_WRITE_PERSISTENT": MISSION_WRITE_PERSISTENT,
	"MISSION_RESET_DEFAULT":    MISSION_RESET_DEFAULT,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e PREFLIGHT_STORAGE_MISSION_ACTION) MarshalText() ([]byte, error) {
	if name, ok := labels_PREFLIGHT_STORAGE_MISSION_ACTION[e]; ok {
		return []byte(name), nil
	}
	return []byte(strconv.Itoa(int(e))), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *PREFLIGHT_STORAGE_MISSION_ACTION) UnmarshalText(text []byte) error {
	if value, ok := values_PREFLIGHT_STORAGE_MISSION_ACTION[string(text)]; ok {
		*e = value
	} else if value, err := strconv.Atoi(string(text)); err == nil {
		*e = PREFLIGHT_STORAGE_MISSION_ACTION(value)
	} else {
		return fmt.Errorf("invalid label '%s'", text)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e PREFLIGHT_STORAGE_MISSION_ACTION) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
