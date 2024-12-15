//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package matrixpilot

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Indicates the ESC connection type.
type ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE

const (
	// Traditional PPM ESC.
	ESC_CONNECTION_TYPE_PPM ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_PPM
	// Serial Bus connected ESC.
	ESC_CONNECTION_TYPE_SERIAL ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_SERIAL
	// One Shot PPM ESC.
	ESC_CONNECTION_TYPE_ONESHOT ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_ONESHOT
	// I2C ESC.
	ESC_CONNECTION_TYPE_I2C ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_I2C
	// CAN-Bus ESC.
	ESC_CONNECTION_TYPE_CAN ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_CAN
	// DShot ESC.
	ESC_CONNECTION_TYPE_DSHOT ESC_CONNECTION_TYPE = common.ESC_CONNECTION_TYPE_DSHOT
)
