//autogenerated:yes
//nolint:revive,misspell,govet,lll
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Update the data in the OPEN_DRONE_ID_SYSTEM message with new location information. This can be sent to update the location information for the operator when no other information in the SYSTEM message has changed. This message allows for efficient operation on radio links which have limited uplink bandwidth while meeting requirements for update frequency of the operator location.
type MessageOpenDroneIdSystemUpdate = common.MessageOpenDroneIdSystemUpdate
