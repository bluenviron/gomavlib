//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package paparazzi

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Standard modes with a well understood meaning across flight stacks and vehicle types.
// For example, most flight stack have the concept of a "return" or "RTL" mode that takes a vehicle to safety, even though the precise mechanics of this mode may differ.
// The modes supported by a flight stack can be queried using AVAILABLE_MODES and set using MAV_CMD_DO_SET_STANDARD_MODE.
// The current mode is streamed in CURRENT_MODE.
// See https://mavlink.io/en/services/standard_modes.html
type MAV_STANDARD_MODE = common.MAV_STANDARD_MODE

const (
	// Non standard mode.
	// This may be used when reporting the mode if the current flight mode is not a standard mode.
	MAV_STANDARD_MODE_NON_STANDARD MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_NON_STANDARD
	// Position mode (manual).
	// Position-controlled and stabilized manual mode.
	// When sticks are released vehicles return to their level-flight orientation and hold both position and altitude against wind and external forces.
	// This mode can only be set by vehicles that can hold a fixed position.
	// Multicopter (MC) vehicles actively brake and hold both position and altitude against wind and external forces.
	// Hybrid MC/FW ("VTOL") vehicles first transition to multicopter mode (if needed) but otherwise behave in the same way as MC vehicles.
	// Fixed-wing (FW) vehicles must not support this mode.
	// Other vehicle types must not support this mode (this may be revisited through the PR process).
	MAV_STANDARD_MODE_POSITION_HOLD MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_POSITION_HOLD
	// Orbit (manual).
	// Position-controlled and stabilized manual mode.
	// The vehicle circles around a fixed setpoint in the horizontal plane at a particular radius, altitude, and direction.
	// Flight stacks may further allow manual control over the setpoint position, radius, direction, speed, and/or altitude of the circle, but this is not mandated.
	// Flight stacks may support the [MAV_CMD_DO_ORBIT](https://mavlink.io/en/messages/common.html#MAV_CMD_DO_ORBIT) for changing the orbit parameters.
	// MC and FW vehicles may support this mode.
	// Hybrid MC/FW ("VTOL") vehicles may support this mode in MC/FW or both modes; if the mode is not supported by the current configuration the vehicle should transition to the supported configuration.
	// Other vehicle types must not support this mode (this may be revisited through the PR process).
	MAV_STANDARD_MODE_ORBIT MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_ORBIT
	// Cruise mode (manual).
	// Position-controlled and stabilized manual mode.
	// When sticks are released vehicles return to their level-flight orientation and hold their original track against wind and external forces.
	// Fixed-wing (FW) vehicles level orientation and maintain current track and altitude against wind and external forces.
	// Hybrid MC/FW ("VTOL") vehicles first transition to FW mode (if needed) but otherwise behave in the same way as MC vehicles.
	// Multicopter (MC) vehicles must not support this mode.
	// Other vehicle types must not support this mode (this may be revisited through the PR process).
	MAV_STANDARD_MODE_CRUISE MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_CRUISE
	// Altitude hold (manual).
	// Altitude-controlled and stabilized manual mode.
	// When sticks are released vehicles return to their level-flight orientation and hold their altitude.
	// MC vehicles continue with existing momentum and may move with wind (or other external forces).
	// FW vehicles continue with current heading, but may be moved off-track by wind.
	// Hybrid MC/FW ("VTOL") vehicles behave according to their current configuration/mode (FW or MC).
	// Other vehicle types must not support this mode (this may be revisited through the PR process).
	MAV_STANDARD_MODE_ALTITUDE_HOLD MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_ALTITUDE_HOLD
	// Safe recovery mode (auto).
	// Automatic mode that takes vehicle to a predefined safe location via a safe flight path, and may also automatically land the vehicle.
	// This mode is more commonly referred to as RTL and/or or Smart RTL.
	// The precise return location, flight path, and landing behaviour depend on vehicle configuration and type.
	// For example, the vehicle might return to the home/launch location, a rally point, or the start of a mission landing, it might follow a direct path, mission path, or breadcrumb path, and land using a mission landing pattern or some other kind of descent.
	MAV_STANDARD_MODE_SAFE_RECOVERY MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_SAFE_RECOVERY
	// Mission mode (automatic).
	// Automatic mode that executes MAVLink missions.
	// Missions are executed from the current waypoint as soon as the mode is enabled.
	MAV_STANDARD_MODE_MISSION MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_MISSION
	// Land mode (auto).
	// Automatic mode that lands the vehicle at the current location.
	// The precise landing behaviour depends on vehicle configuration and type.
	MAV_STANDARD_MODE_LAND MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_LAND
	// Takeoff mode (auto).
	// Automatic takeoff mode.
	// The precise takeoff behaviour depends on vehicle configuration and type.
	MAV_STANDARD_MODE_TAKEOFF MAV_STANDARD_MODE = common.MAV_STANDARD_MODE_TAKEOFF
)
