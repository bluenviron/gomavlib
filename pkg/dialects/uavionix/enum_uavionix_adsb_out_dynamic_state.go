//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package uavionix

import (
	"fmt"
	"strconv"
	"strings"
)

// State flags for ADS-B transponder dynamic report
type UAVIONIX_ADSB_OUT_DYNAMIC_STATE uint32

const (
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE        UAVIONIX_ADSB_OUT_DYNAMIC_STATE = 1
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED    UAVIONIX_ADSB_OUT_DYNAMIC_STATE = 2
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED UAVIONIX_ADSB_OUT_DYNAMIC_STATE = 4
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND            UAVIONIX_ADSB_OUT_DYNAMIC_STATE = 8
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT                UAVIONIX_ADSB_OUT_DYNAMIC_STATE = 16
)

var labels_UAVIONIX_ADSB_OUT_DYNAMIC_STATE = map[UAVIONIX_ADSB_OUT_DYNAMIC_STATE]string{
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE:        "UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE",
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED:    "UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED",
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED: "UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED",
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND:            "UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND",
	UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT:                "UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT",
}

var values_UAVIONIX_ADSB_OUT_DYNAMIC_STATE = map[string]UAVIONIX_ADSB_OUT_DYNAMIC_STATE{
	"UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE":        UAVIONIX_ADSB_OUT_DYNAMIC_STATE_INTENT_CHANGE,
	"UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED":    UAVIONIX_ADSB_OUT_DYNAMIC_STATE_AUTOPILOT_ENABLED,
	"UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED": UAVIONIX_ADSB_OUT_DYNAMIC_STATE_NICBARO_CROSSCHECKED,
	"UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND":            UAVIONIX_ADSB_OUT_DYNAMIC_STATE_ON_GROUND,
	"UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT":                UAVIONIX_ADSB_OUT_DYNAMIC_STATE_IDENT,
}

// MarshalText implements the encoding.TextMarshaler interface.
func (e UAVIONIX_ADSB_OUT_DYNAMIC_STATE) MarshalText() ([]byte, error) {
	if e == 0 {
		return []byte("0"), nil
	}
	var names []string
	for i := 0; i < 5; i++ {
		mask := UAVIONIX_ADSB_OUT_DYNAMIC_STATE(1 << i)
		if e&mask == mask {
			names = append(names, labels_UAVIONIX_ADSB_OUT_DYNAMIC_STATE[mask])
		}
	}
	return []byte(strings.Join(names, " | ")), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (e *UAVIONIX_ADSB_OUT_DYNAMIC_STATE) UnmarshalText(text []byte) error {
	labels := strings.Split(string(text), " | ")
	var mask UAVIONIX_ADSB_OUT_DYNAMIC_STATE
	for _, label := range labels {
		if value, ok := values_UAVIONIX_ADSB_OUT_DYNAMIC_STATE[label]; ok {
			mask |= value
		} else if value, err := strconv.Atoi(label); err == nil {
			mask |= UAVIONIX_ADSB_OUT_DYNAMIC_STATE(value)
		} else {
			return fmt.Errorf("invalid label '%s'", label)
		}
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (e UAVIONIX_ADSB_OUT_DYNAMIC_STATE) String() string {
	val, _ := e.MarshalText()
	return string(val)
}
