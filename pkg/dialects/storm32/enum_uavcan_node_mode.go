//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package storm32

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Generalized UAVCAN node mode
type UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE

const (
	// The node is performing its primary functions.
	UAVCAN_NODE_MODE_OPERATIONAL UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE_OPERATIONAL
	// The node is initializing; this mode is entered immediately after startup.
	UAVCAN_NODE_MODE_INITIALIZATION UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE_INITIALIZATION
	// The node is under maintenance.
	UAVCAN_NODE_MODE_MAINTENANCE UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE_MAINTENANCE
	// The node is in the process of updating its software.
	UAVCAN_NODE_MODE_SOFTWARE_UPDATE UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE_SOFTWARE_UPDATE
	// The node is no longer available online.
	UAVCAN_NODE_MODE_OFFLINE UAVCAN_NODE_MODE = common.UAVCAN_NODE_MODE_OFFLINE
)
