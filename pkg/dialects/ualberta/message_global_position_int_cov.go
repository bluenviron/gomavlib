//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ualberta

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// The filtered global position (e.g. fused GPS and accelerometers). The position is in GPS-frame (right-handed, Z-up). It  is designed as scaled integer message since the resolution of float is not sufficient. NOTE: This message is intended for onboard networks / companion computers and higher-bandwidth links and optimized for accuracy and completeness. Please use the GLOBAL_POSITION_INT message for a minimal subset.
type MessageGlobalPositionIntCov = common.MessageGlobalPositionIntCov
