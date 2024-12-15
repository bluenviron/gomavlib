//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Battery information. Updates GCS with flight controller battery status. Smart batteries also use this message, but may additionally send BATTERY_INFO.
type MessageBatteryStatus = common.MessageBatteryStatus
