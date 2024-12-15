//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Enumeration for battery charge states.
type MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE

const (
	// Low battery state is not provided
	MAV_BATTERY_CHARGE_STATE_UNDEFINED MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_UNDEFINED
	// Battery is not in low state. Normal operation.
	MAV_BATTERY_CHARGE_STATE_OK MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_OK
	// Battery state is low, warn and monitor close.
	MAV_BATTERY_CHARGE_STATE_LOW MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_LOW
	// Battery state is critical, return or abort immediately.
	MAV_BATTERY_CHARGE_STATE_CRITICAL MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_CRITICAL
	// Battery state is too low for ordinary abort sequence. Perform fastest possible emergency stop to prevent damage.
	MAV_BATTERY_CHARGE_STATE_EMERGENCY MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_EMERGENCY
	// Battery failed, damage unavoidable. Possible causes (faults) are listed in MAV_BATTERY_FAULT.
	MAV_BATTERY_CHARGE_STATE_FAILED MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_FAILED
	// Battery is diagnosed to be defective or an error occurred, usage is discouraged / prohibited. Possible causes (faults) are listed in MAV_BATTERY_FAULT.
	MAV_BATTERY_CHARGE_STATE_UNHEALTHY MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_UNHEALTHY
	// Battery is charging.
	MAV_BATTERY_CHARGE_STATE_CHARGING MAV_BATTERY_CHARGE_STATE = common.MAV_BATTERY_CHARGE_STATE_CHARGING
)
