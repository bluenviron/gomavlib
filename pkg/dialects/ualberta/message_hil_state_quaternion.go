//autogenerated:yes
//nolint:revive,misspell,govet,lll
package ualberta

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Sent from simulation to autopilot, avoids in contrast to HIL_STATE singularities. This packet is useful for high throughput applications such as hardware in the loop simulations.
type MessageHilStateQuaternion = common.MessageHilStateQuaternion
