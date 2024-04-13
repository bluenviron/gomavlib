//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package ualberta

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Power supply status flags (bitmask)
type MAV_POWER_STATUS = common.MAV_POWER_STATUS

const (
	// main brick power supply valid
	MAV_POWER_STATUS_BRICK_VALID MAV_POWER_STATUS = common.MAV_POWER_STATUS_BRICK_VALID
	// main servo power supply valid for FMU
	MAV_POWER_STATUS_SERVO_VALID MAV_POWER_STATUS = common.MAV_POWER_STATUS_SERVO_VALID
	// USB power is connected
	MAV_POWER_STATUS_USB_CONNECTED MAV_POWER_STATUS = common.MAV_POWER_STATUS_USB_CONNECTED
	// peripheral supply is in over-current state
	MAV_POWER_STATUS_PERIPH_OVERCURRENT MAV_POWER_STATUS = common.MAV_POWER_STATUS_PERIPH_OVERCURRENT
	// hi-power peripheral supply is in over-current state
	MAV_POWER_STATUS_PERIPH_HIPOWER_OVERCURRENT MAV_POWER_STATUS = common.MAV_POWER_STATUS_PERIPH_HIPOWER_OVERCURRENT
	// Power status has changed since boot
	MAV_POWER_STATUS_CHANGED MAV_POWER_STATUS = common.MAV_POWER_STATUS_CHANGED
)
