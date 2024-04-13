//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

type FENCE_TYPE = common.FENCE_TYPE

const (
	// All fence types
	FENCE_TYPE_ALL FENCE_TYPE = common.FENCE_TYPE_ALL
	// Maximum altitude fence
	FENCE_TYPE_ALT_MAX FENCE_TYPE = common.FENCE_TYPE_ALT_MAX
	// Circle fence
	FENCE_TYPE_CIRCLE FENCE_TYPE = common.FENCE_TYPE_CIRCLE
	// Polygon fence
	FENCE_TYPE_POLYGON FENCE_TYPE = common.FENCE_TYPE_POLYGON
	// Minimum altitude fence
	FENCE_TYPE_ALT_MIN FENCE_TYPE = common.FENCE_TYPE_ALT_MIN
)
