//autogenerated:yes
//nolint:revive,misspell,govet,lll
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Request one or more events to be (re-)sent. If first_sequence==last_sequence, only a single event is requested. Note that first_sequence can be larger than last_sequence (because the sequence number can wrap). Each sequence will trigger an EVENT or EVENT_ERROR response.
type MessageRequestEvent = common.MessageRequestEvent
