//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package common

import (
	"fmt"
	"strconv"
)

// Flags to indicate the type of storage.
type STORAGE_TYPE uint32

const (
	// Storage type is not known.
	STORAGE_TYPE_UNKNOWN STORAGE_TYPE = 0
	// Storage type is USB device.
	STORAGE_TYPE_USB_STICK STORAGE_TYPE = 1
	// Storage type is SD card.
	STORAGE_TYPE_SD STORAGE_TYPE = 2
	// Storage type is microSD card.
	STORAGE_TYPE_MICROSD STORAGE_TYPE = 3
	// Storage type is CFast.
	STORAGE_TYPE_CF STORAGE_TYPE = 4
	// Storage type is CFexpress.
	STORAGE_TYPE_CFE STORAGE_TYPE = 5
	// Storage type is XQD.
	STORAGE_TYPE_XQD STORAGE_TYPE = 6
	// Storage type is HD mass storage type.
	STORAGE_TYPE_HD STORAGE_TYPE = 7
	// Storage type is other, not listed type.
	STORAGE_TYPE_OTHER STORAGE_TYPE = 254
)

var labels_STORAGE_TYPE = map[STORAGE_TYPE]string{
	STORAGE_TYPE_UNKNOWN:   "STORAGE_TYPE_UNKNOWN",
	STORAGE_TYPE_USB_STICK: "STORAGE_TYPE_USB_STICK",
	STORAGE_TYPE_SD:        "STORAGE_TYPE_SD",
	STORAGE_TYPE_MICROSD:   "STORAGE_TYPE_MICROSD",
	STORAGE_TYPE_CF:        "STORAGE_TYPE_CF",
	STORAGE_TYPE_CFE:       "STORAGE_TYPE_CFE",
	STORAGE_TYPE_XQD:       "STORAGE_TYPE_XQD",
	STORAGE_TYPE_HD:        "STORAGE_TYPE_HD",
	STORAGE_TYPE_OTHER:     "STORAGE_TYPE_OTHER",
}

var values_STORAGE_TYPE = map[string]STORAGE_TYPE{
	"STORAGE_TYPE_UNKNOWN":   STORAGE_TYPE_UNKNOWN,
	"STORAGE_TYPE_USB_STICK": STORAGE_TYPE_USB_STICK,
	"STORAGE_TYPE_SD":        STORAGE_TYPE_SD,
	"STORAGE_TYPE_MICROSD":   STORAGE_TYPE_MICROSD,
	"STORAGE_TYPE_CF":        STORAGE_TYPE_CF,
	"STORAGE_TYPE_CFE":       STORAGE_TYPE_CFE,
	"STORAGE_TYPE_XQD":       STORAGE_TYPE_XQD,
	"STORAGE_TYPE_HD":        STORAGE_TYPE_HD,
	"STORAGE_TYPE_OTHER":     STORAGE_TYPE_OTHER,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e STORAGE_TYPE) MarshalText() ([]byte, error) {
	name, ok := labels_STORAGE_TYPE[e]
	if !ok {
		return nil, fmt.Errorf("invalid value %d", e)
	}
	return []byte(name), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *STORAGE_TYPE) UnmarshalText(text []byte) error {
	value, ok := values_STORAGE_TYPE[string(text)]
	if !ok {
		return fmt.Errorf("invalid label '%s'", text)
	}
	*e = value
	return nil
}

// String implements the fmt.Stringer interface.
func (e STORAGE_TYPE) String() string {
	name, ok := labels_STORAGE_TYPE[e]
	if !ok {
		return strconv.Itoa(int(e))
	}
	return name
}
