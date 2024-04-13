//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package asluav

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Flags to indicate the type of storage.
type STORAGE_TYPE = common.STORAGE_TYPE

const (
	// Storage type is not known.
	STORAGE_TYPE_UNKNOWN STORAGE_TYPE = common.STORAGE_TYPE_UNKNOWN
	// Storage type is USB device.
	STORAGE_TYPE_USB_STICK STORAGE_TYPE = common.STORAGE_TYPE_USB_STICK
	// Storage type is SD card.
	STORAGE_TYPE_SD STORAGE_TYPE = common.STORAGE_TYPE_SD
	// Storage type is microSD card.
	STORAGE_TYPE_MICROSD STORAGE_TYPE = common.STORAGE_TYPE_MICROSD
	// Storage type is CFast.
	STORAGE_TYPE_CF STORAGE_TYPE = common.STORAGE_TYPE_CF
	// Storage type is CFexpress.
	STORAGE_TYPE_CFE STORAGE_TYPE = common.STORAGE_TYPE_CFE
	// Storage type is XQD.
	STORAGE_TYPE_XQD STORAGE_TYPE = common.STORAGE_TYPE_XQD
	// Storage type is HD mass storage type.
	STORAGE_TYPE_HD STORAGE_TYPE = common.STORAGE_TYPE_HD
	// Storage type is other, not listed type.
	STORAGE_TYPE_OTHER STORAGE_TYPE = common.STORAGE_TYPE_OTHER
)
