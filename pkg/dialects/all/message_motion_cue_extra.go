//autogenerated:yes
//nolint:revive,misspell,govet,lll
package all

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/marsh"
)

// These values are an extra cue that should be added to accelerations and rotations etc. resulting from aircraft state, with the resulting cue being the sum of the latest aircraft and extra values. An example use case would be a cockpit shaker.
type MessageMotionCueExtra = marsh.MessageMotionCueExtra
