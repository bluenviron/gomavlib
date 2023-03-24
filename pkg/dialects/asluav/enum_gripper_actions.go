//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package asluav

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Gripper actions.
type GRIPPER_ACTIONS = common.GRIPPER_ACTIONS

const (
	// Gripper release cargo.
	GRIPPER_ACTION_RELEASE GRIPPER_ACTIONS = common.GRIPPER_ACTION_RELEASE
	// Gripper grab onto cargo.
	GRIPPER_ACTION_GRAB GRIPPER_ACTIONS = common.GRIPPER_ACTION_GRAB
)
