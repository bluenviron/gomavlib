//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package uavionix

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// These values define the type of firmware release.  These values indicate the first version or release of this type.  For example the first alpha release would be 64, the second would be 65.
type FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE

const (
	// development release
	FIRMWARE_VERSION_TYPE_DEV FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE_DEV
	// alpha release
	FIRMWARE_VERSION_TYPE_ALPHA FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE_ALPHA
	// beta release
	FIRMWARE_VERSION_TYPE_BETA FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE_BETA
	// release candidate
	FIRMWARE_VERSION_TYPE_RC FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE_RC
	// official stable release
	FIRMWARE_VERSION_TYPE_OFFICIAL FIRMWARE_VERSION_TYPE = common.FIRMWARE_VERSION_TYPE_OFFICIAL
)
