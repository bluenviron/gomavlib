//autogenerated:yes
//nolint:revive,misspell,govet,lll
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// High level message to control a gimbal manually. The angles or angular rates are unitless; the actual rates will depend on internal gimbal manager settings/configuration (e.g. set by parameters). This message is to be sent to the gimbal manager (e.g. from a ground station). Angles and rates can be set to NaN according to use case.
type MessageGimbalManagerSetManualControl = common.MessageGimbalManagerSetManualControl
