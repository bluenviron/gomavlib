//autogenerated:yes
//nolint:revive,misspell,govet,lll
package all

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Send a command with up to seven parameters to the MAV. COMMAND_INT is generally preferred when sending MAV_CMD commands that include positional information; it offers higher precision and allows the MAV_FRAME to be specified (which may otherwise be ambiguous, particularly for altitude). The command microservice is documented at https://mavlink.io/en/services/command.html
type MessageCommandLong = common.MessageCommandLong
