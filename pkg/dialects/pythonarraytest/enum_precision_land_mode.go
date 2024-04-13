//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package pythonarraytest

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Precision land modes (used in MAV_CMD_NAV_LAND).
type PRECISION_LAND_MODE = common.PRECISION_LAND_MODE

const (
	// Normal (non-precision) landing.
	PRECISION_LAND_MODE_DISABLED PRECISION_LAND_MODE = common.PRECISION_LAND_MODE_DISABLED
	// Use precision landing if beacon detected when land command accepted, otherwise land normally.
	PRECISION_LAND_MODE_OPPORTUNISTIC PRECISION_LAND_MODE = common.PRECISION_LAND_MODE_OPPORTUNISTIC
	// Use precision landing, searching for beacon if not found when land command accepted (land normally if beacon cannot be found).
	PRECISION_LAND_MODE_REQUIRED PRECISION_LAND_MODE = common.PRECISION_LAND_MODE_REQUIRED
)
