//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// Possible remote log data block statuses.
type MAV_REMOTE_LOG_DATA_BLOCK_STATUSES = ardupilotmega.MAV_REMOTE_LOG_DATA_BLOCK_STATUSES

const (
	// This block has NOT been received.
	MAV_REMOTE_LOG_DATA_BLOCK_NACK MAV_REMOTE_LOG_DATA_BLOCK_STATUSES = ardupilotmega.MAV_REMOTE_LOG_DATA_BLOCK_NACK
	// This block has been received.
	MAV_REMOTE_LOG_DATA_BLOCK_ACK MAV_REMOTE_LOG_DATA_BLOCK_STATUSES = ardupilotmega.MAV_REMOTE_LOG_DATA_BLOCK_ACK
)
