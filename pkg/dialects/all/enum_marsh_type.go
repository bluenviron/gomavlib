//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package all

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/marsh"
)

// Component types for different nodes of the simulator network (flight model, controls, visualisation etc.). Components will always receive messages from the Manager relevant for their type. Only the first component in a network with a given system ID and type will have its messages forwarded by the Manager, all other ones will only be treated as output (will be shadowed). This enum is an extension of MAV_TYPE documented at https://mavlink.io/en/messages/minimal.html#MAV_TYPE
type MARSH_TYPE = marsh.MARSH_TYPE

const (
	// The simulation manager responsible for routing packets between different nodes. Typically MARSH Manager, see https://marsh-sim.github.io/manager.html
	MARSH_TYPE_MANAGER MARSH_TYPE = marsh.MARSH_TYPE_MANAGER
	// Component simulating flight dynamics of the aircraft.
	MARSH_TYPE_FLIGHT_MODEL MARSH_TYPE = marsh.MARSH_TYPE_FLIGHT_MODEL
	// Component providing pilot control inputs.
	MARSH_TYPE_CONTROLS MARSH_TYPE = marsh.MARSH_TYPE_CONTROLS
	// Component showing the visual situation to the pilot.
	MARSH_TYPE_VISUALISATION MARSH_TYPE = marsh.MARSH_TYPE_VISUALISATION
	// Component implementing pilot instrument panel.
	MARSH_TYPE_INSTRUMENTS MARSH_TYPE = marsh.MARSH_TYPE_INSTRUMENTS
	// Component that moves the entire cockpit for motion cueing.
	MARSH_TYPE_MOTION_PLATFORM MARSH_TYPE = marsh.MARSH_TYPE_MOTION_PLATFORM
	// Component for in-seat motion cueing.
	MARSH_TYPE_GSEAT MARSH_TYPE = marsh.MARSH_TYPE_GSEAT
	// Component providing gaze data of pilot eyes.
	MARSH_TYPE_EYE_TRACKER MARSH_TYPE = marsh.MARSH_TYPE_EYE_TRACKER
	// Component measuring and actuating forces on pilot control inputs.
	MARSH_TYPE_CONTROL_LOADING MARSH_TYPE = marsh.MARSH_TYPE_CONTROL_LOADING
	// Component providing vibrations for system identification, road rumble, gusts, etc.
	MARSH_TYPE_VIBRATION_SOURCE MARSH_TYPE = marsh.MARSH_TYPE_VIBRATION_SOURCE
	// Component providing target for the pilot to follow like controls positions, aircraft state, ILS path etc.
	MARSH_TYPE_PILOT_TARGET MARSH_TYPE = marsh.MARSH_TYPE_PILOT_TARGET
	// Principal component controlling the main scenario of a given test, (unlike lower level MARSH_TYPE_PILOT_TARGET or MARSH_TYPE_MANAGER for overall communication).
	MARSH_TYPE_EXPERIMENT_DIRECTOR MARSH_TYPE = marsh.MARSH_TYPE_EXPERIMENT_DIRECTOR
)
