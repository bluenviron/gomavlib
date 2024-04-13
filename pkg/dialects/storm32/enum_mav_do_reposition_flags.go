//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Bitmap of options for the MAV_CMD_DO_REPOSITION
type MAV_DO_REPOSITION_FLAGS = common.MAV_DO_REPOSITION_FLAGS

const (
	// The aircraft should immediately transition into guided. This should not be set for follow me applications
	MAV_DO_REPOSITION_FLAGS_CHANGE_MODE MAV_DO_REPOSITION_FLAGS = common.MAV_DO_REPOSITION_FLAGS_CHANGE_MODE
)
