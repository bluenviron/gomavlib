//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ardupilotmega

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// This message is sent to the MAV to write a partial list. If start index == end index, only one item will be transmitted / updated. If the start index is NOT 0 and above the current list size, this request should be REJECTED!
type MessageMissionWritePartialList = common.MessageMissionWritePartialList
