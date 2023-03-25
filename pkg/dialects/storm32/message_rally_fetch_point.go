//autogenerated:yes
//nolint:revive,misspell,govet,lll
package storm32

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
)

// Request a current rally point from MAV. MAV should respond with a RALLY_POINT message. MAV should not respond if the request is invalid.
type MessageRallyFetchPoint = ardupilotmega.MessageRallyFetchPoint
