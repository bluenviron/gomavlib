//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/minimal"
)

// The heartbeat message shows that a system or component is present and responding. The type and autopilot fields (along with the message component id), allow the receiving system to treat further messages from this system appropriately (e.g. by laying out the user interface based on the autopilot). This microservice is documented at https://mavlink.io/en/services/heartbeat.html
type MessageHeartbeat = minimal.MessageHeartbeat
