//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/development"
)

// Mode properties.
type MAV_MODE_PROPERTY = development.MAV_MODE_PROPERTY

const (
	// If set, this mode is an advanced mode.
	// For example a rate-controlled manual mode might be advanced, whereas a position-controlled manual mode is not.
	// A GCS can optionally use this flag to configure the UI for its intended users.
	MAV_MODE_PROPERTY_ADVANCED MAV_MODE_PROPERTY = development.MAV_MODE_PROPERTY_ADVANCED
	// If set, this mode should not be added to the list of selectable modes.
	// The mode might still be selected by the FC directly (for example as part of a failsafe).
	MAV_MODE_PROPERTY_NOT_USER_SELECTABLE MAV_MODE_PROPERTY = development.MAV_MODE_PROPERTY_NOT_USER_SELECTABLE
)
