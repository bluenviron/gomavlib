//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

type GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAGS

const (
	// ignore altitude field
	GPS_INPUT_IGNORE_FLAG_ALT GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_ALT
	// ignore hdop field
	GPS_INPUT_IGNORE_FLAG_HDOP GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_HDOP
	// ignore vdop field
	GPS_INPUT_IGNORE_FLAG_VDOP GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_VDOP
	// ignore horizontal velocity field (vn and ve)
	GPS_INPUT_IGNORE_FLAG_VEL_HORIZ GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_VEL_HORIZ
	// ignore vertical velocity field (vd)
	GPS_INPUT_IGNORE_FLAG_VEL_VERT GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_VEL_VERT
	// ignore speed accuracy field
	GPS_INPUT_IGNORE_FLAG_SPEED_ACCURACY GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_SPEED_ACCURACY
	// ignore horizontal accuracy field
	GPS_INPUT_IGNORE_FLAG_HORIZONTAL_ACCURACY GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_HORIZONTAL_ACCURACY
	// ignore vertical accuracy field
	GPS_INPUT_IGNORE_FLAG_VERTICAL_ACCURACY GPS_INPUT_IGNORE_FLAGS = common.GPS_INPUT_IGNORE_FLAG_VERTICAL_ACCURACY
)
