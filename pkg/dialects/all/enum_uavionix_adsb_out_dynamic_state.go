//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/uavionix"
)

// State flags for ADS-B transponder dynamic report
type UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE

const (
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE        UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED    UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND            UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT                UAVIONIX_ADSB_OUT_DYNAMIC_STATE = uavionix.UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT
)
