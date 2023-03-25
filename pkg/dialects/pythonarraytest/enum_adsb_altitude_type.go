//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package pythonarraytest

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Enumeration of the ADSB altimeter types
type ADSB_ALTITUDE_TYPE = common.ADSB_ALTITUDE_TYPE

const (
	// Altitude reported from a Baro source using QNH reference
	ADSB_ALTITUDE_TYPE_PRESSURE_QNH ADSB_ALTITUDE_TYPE = common.ADSB_ALTITUDE_TYPE_PRESSURE_QNH
	// Altitude reported from a GNSS source
	ADSB_ALTITUDE_TYPE_GEOMETRIC ADSB_ALTITUDE_TYPE = common.ADSB_ALTITUDE_TYPE_GEOMETRIC
)
