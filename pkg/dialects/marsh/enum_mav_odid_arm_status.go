//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package marsh

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

type MAV_ODID_ARM_STATUS = common.MAV_ODID_ARM_STATUS

const (
	// Passing arming checks.
	MAV_ODID_ARM_STATUS_GOOD_TO_ARM MAV_ODID_ARM_STATUS = common.MAV_ODID_ARM_STATUS_GOOD_TO_ARM
	// Generic arming failure, see error string for details.
	MAV_ODID_ARM_STATUS_PRE_ARM_FAIL_GENERIC MAV_ODID_ARM_STATUS = common.MAV_ODID_ARM_STATUS_PRE_ARM_FAIL_GENERIC
)
