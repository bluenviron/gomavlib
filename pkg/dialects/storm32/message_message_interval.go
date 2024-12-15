//autogenerated:yes
//nolint:revive,misspell,govet,lll
package storm32

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// The interval between messages for a particular MAVLink message ID.
// This message is sent in response to the MAV_CMD_REQUEST_MESSAGE command with param1=244 (this message) and param2=message_id (the id of the message for which the interval is required).
// It may also be sent in response to MAV_CMD_GET_MESSAGE_INTERVAL.
// This interface replaces DATA_STREAM.
type MessageMessageInterval = common.MessageMessageInterval
