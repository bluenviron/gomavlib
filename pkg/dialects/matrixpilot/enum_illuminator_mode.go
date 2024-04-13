//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package matrixpilot

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Modes of illuminator
type ILLUMINATOR_MODE = common.ILLUMINATOR_MODE

const (
	// Illuminator mode is not specified/unknown
	ILLUMINATOR_MODE_UNKNOWN ILLUMINATOR_MODE = common.ILLUMINATOR_MODE_UNKNOWN
	// Illuminator behavior is controlled by MAV_CMD_DO_ILLUMINATOR_CONFIGURE settings
	ILLUMINATOR_MODE_INTERNAL_CONTROL ILLUMINATOR_MODE = common.ILLUMINATOR_MODE_INTERNAL_CONTROL
	// Illuminator behavior is controlled by external factors: e.g. an external hardware signal
	ILLUMINATOR_MODE_EXTERNAL_SYNC ILLUMINATOR_MODE = common.ILLUMINATOR_MODE_EXTERNAL_SYNC
)
