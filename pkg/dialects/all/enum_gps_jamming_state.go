//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/development"
)

// Signal jamming state in a GPS receiver.
type GPS_JAMMING_STATE = development.GPS_JAMMING_STATE

const (
	// The GPS receiver does not provide GPS signal jamming info.
	GPS_JAMMING_STATE_UNKNOWN GPS_JAMMING_STATE = development.GPS_JAMMING_STATE_UNKNOWN
	// The GPS receiver detected no signal jamming.
	GPS_JAMMING_STATE_OK GPS_JAMMING_STATE = development.GPS_JAMMING_STATE_OK
	// The GPS receiver detected and mitigated signal jamming.
	GPS_JAMMING_STATE_MITIGATED GPS_JAMMING_STATE = development.GPS_JAMMING_STATE_MITIGATED
	// The GPS receiver detected signal jamming.
	GPS_JAMMING_STATE_DETECTED GPS_JAMMING_STATE = development.GPS_JAMMING_STATE_DETECTED
)
