//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Type of landing target
type LANDING_TARGET_TYPE = common.LANDING_TARGET_TYPE

const (
	// Landing target signaled by light beacon (ex: IR-LOCK)
	LANDING_TARGET_TYPE_LIGHT_BEACON LANDING_TARGET_TYPE = common.LANDING_TARGET_TYPE_LIGHT_BEACON
	// Landing target signaled by radio beacon (ex: ILS, NDB)
	LANDING_TARGET_TYPE_RADIO_BEACON LANDING_TARGET_TYPE = common.LANDING_TARGET_TYPE_RADIO_BEACON
	// Landing target represented by a fiducial marker (ex: ARTag)
	LANDING_TARGET_TYPE_VISION_FIDUCIAL LANDING_TARGET_TYPE = common.LANDING_TARGET_TYPE_VISION_FIDUCIAL
	// Landing target represented by a pre-defined visual shape/feature (ex: X-marker, H-marker, square)
	LANDING_TARGET_TYPE_VISION_OTHER LANDING_TARGET_TYPE = common.LANDING_TARGET_TYPE_VISION_OTHER
)
