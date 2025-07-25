//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package matrixpilot

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/minimal"
)

// MAVLINK component type reported in HEARTBEAT message. Flight controllers must report the type of the vehicle on which they are mounted (e.g. MAV_TYPE_OCTOROTOR). All other components must report a value appropriate for their type (e.g. a camera must use MAV_TYPE_CAMERA).
type MAV_TYPE = minimal.MAV_TYPE

const (
	// Generic micro air vehicle
	MAV_TYPE_GENERIC MAV_TYPE = minimal.MAV_TYPE_GENERIC
	// Fixed wing aircraft.
	MAV_TYPE_FIXED_WING MAV_TYPE = minimal.MAV_TYPE_FIXED_WING
	// Quadrotor
	MAV_TYPE_QUADROTOR MAV_TYPE = minimal.MAV_TYPE_QUADROTOR
	// Coaxial helicopter
	MAV_TYPE_COAXIAL MAV_TYPE = minimal.MAV_TYPE_COAXIAL
	// Normal helicopter with tail rotor.
	MAV_TYPE_HELICOPTER MAV_TYPE = minimal.MAV_TYPE_HELICOPTER
	// Ground installation
	MAV_TYPE_ANTENNA_TRACKER MAV_TYPE = minimal.MAV_TYPE_ANTENNA_TRACKER
	// Operator control unit / ground control station
	MAV_TYPE_GCS MAV_TYPE = minimal.MAV_TYPE_GCS
	// Airship, controlled
	MAV_TYPE_AIRSHIP MAV_TYPE = minimal.MAV_TYPE_AIRSHIP
	// Free balloon, uncontrolled
	MAV_TYPE_FREE_BALLOON MAV_TYPE = minimal.MAV_TYPE_FREE_BALLOON
	// Rocket
	MAV_TYPE_ROCKET MAV_TYPE = minimal.MAV_TYPE_ROCKET
	// Ground rover
	MAV_TYPE_GROUND_ROVER MAV_TYPE = minimal.MAV_TYPE_GROUND_ROVER
	// Surface vessel, boat, ship
	MAV_TYPE_SURFACE_BOAT MAV_TYPE = minimal.MAV_TYPE_SURFACE_BOAT
	// Submarine
	MAV_TYPE_SUBMARINE MAV_TYPE = minimal.MAV_TYPE_SUBMARINE
	// Hexarotor
	MAV_TYPE_HEXAROTOR MAV_TYPE = minimal.MAV_TYPE_HEXAROTOR
	// Octorotor
	MAV_TYPE_OCTOROTOR MAV_TYPE = minimal.MAV_TYPE_OCTOROTOR
	// Tricopter
	MAV_TYPE_TRICOPTER MAV_TYPE = minimal.MAV_TYPE_TRICOPTER
	// Flapping wing
	MAV_TYPE_FLAPPING_WING MAV_TYPE = minimal.MAV_TYPE_FLAPPING_WING
	// Kite
	MAV_TYPE_KITE MAV_TYPE = minimal.MAV_TYPE_KITE
	// Onboard companion controller
	MAV_TYPE_ONBOARD_CONTROLLER MAV_TYPE = minimal.MAV_TYPE_ONBOARD_CONTROLLER
	// Two-rotor Tailsitter VTOL that additionally uses control surfaces in vertical operation. Note, value previously named MAV_TYPE_VTOL_DUOROTOR.
	MAV_TYPE_VTOL_TAILSITTER_DUOROTOR MAV_TYPE = minimal.MAV_TYPE_VTOL_TAILSITTER_DUOROTOR
	// Quad-rotor Tailsitter VTOL using a V-shaped quad config in vertical operation. Note: value previously named MAV_TYPE_VTOL_QUADROTOR.
	MAV_TYPE_VTOL_TAILSITTER_QUADROTOR MAV_TYPE = minimal.MAV_TYPE_VTOL_TAILSITTER_QUADROTOR
	// Tiltrotor VTOL. Fuselage and wings stay (nominally) horizontal in all flight phases. It able to tilt (some) rotors to provide thrust in cruise flight.
	MAV_TYPE_VTOL_TILTROTOR MAV_TYPE = minimal.MAV_TYPE_VTOL_TILTROTOR
	// VTOL with separate fixed rotors for hover and cruise flight. Fuselage and wings stay (nominally) horizontal in all flight phases.
	MAV_TYPE_VTOL_FIXEDROTOR MAV_TYPE = minimal.MAV_TYPE_VTOL_FIXEDROTOR
	// Tailsitter VTOL. Fuselage and wings orientation changes depending on flight phase: vertical for hover, horizontal for cruise. Use more specific VTOL MAV_TYPE_VTOL_TAILSITTER_DUOROTOR or MAV_TYPE_VTOL_TAILSITTER_QUADROTOR if appropriate.
	MAV_TYPE_VTOL_TAILSITTER MAV_TYPE = minimal.MAV_TYPE_VTOL_TAILSITTER
	// Tiltwing VTOL. Fuselage stays horizontal in all flight phases. The whole wing, along with any attached engine, can tilt between vertical and horizontal mode.
	MAV_TYPE_VTOL_TILTWING MAV_TYPE = minimal.MAV_TYPE_VTOL_TILTWING
	// VTOL reserved 5
	MAV_TYPE_VTOL_RESERVED5 MAV_TYPE = minimal.MAV_TYPE_VTOL_RESERVED5
	// Gimbal
	MAV_TYPE_GIMBAL MAV_TYPE = minimal.MAV_TYPE_GIMBAL
	// ADSB system
	MAV_TYPE_ADSB MAV_TYPE = minimal.MAV_TYPE_ADSB
	// Steerable, nonrigid airfoil
	MAV_TYPE_PARAFOIL MAV_TYPE = minimal.MAV_TYPE_PARAFOIL
	// Dodecarotor
	MAV_TYPE_DODECAROTOR MAV_TYPE = minimal.MAV_TYPE_DODECAROTOR
	// Camera
	MAV_TYPE_CAMERA MAV_TYPE = minimal.MAV_TYPE_CAMERA
	// Charging station
	MAV_TYPE_CHARGING_STATION MAV_TYPE = minimal.MAV_TYPE_CHARGING_STATION
	// FLARM collision avoidance system
	MAV_TYPE_FLARM MAV_TYPE = minimal.MAV_TYPE_FLARM
	// Servo
	MAV_TYPE_SERVO MAV_TYPE = minimal.MAV_TYPE_SERVO
	// Open Drone ID. See https://mavlink.io/en/services/opendroneid.html.
	MAV_TYPE_ODID MAV_TYPE = minimal.MAV_TYPE_ODID
	// Decarotor
	MAV_TYPE_DECAROTOR MAV_TYPE = minimal.MAV_TYPE_DECAROTOR
	// Battery
	MAV_TYPE_BATTERY MAV_TYPE = minimal.MAV_TYPE_BATTERY
	// Parachute
	MAV_TYPE_PARACHUTE MAV_TYPE = minimal.MAV_TYPE_PARACHUTE
	// Log
	MAV_TYPE_LOG MAV_TYPE = minimal.MAV_TYPE_LOG
	// OSD
	MAV_TYPE_OSD MAV_TYPE = minimal.MAV_TYPE_OSD
	// IMU
	MAV_TYPE_IMU MAV_TYPE = minimal.MAV_TYPE_IMU
	// GPS
	MAV_TYPE_GPS MAV_TYPE = minimal.MAV_TYPE_GPS
	// Winch
	MAV_TYPE_WINCH MAV_TYPE = minimal.MAV_TYPE_WINCH
	// Generic multirotor that does not fit into a specific type or whose type is unknown
	MAV_TYPE_GENERIC_MULTIROTOR MAV_TYPE = minimal.MAV_TYPE_GENERIC_MULTIROTOR
	// Illuminator. An illuminator is a light source that is used for lighting up dark areas external to the system: e.g. a torch or searchlight (as opposed to a light source for illuminating the system itself, e.g. an indicator light).
	MAV_TYPE_ILLUMINATOR MAV_TYPE = minimal.MAV_TYPE_ILLUMINATOR
	// Orbiter spacecraft. Includes satellites orbiting terrestrial and extra-terrestrial bodies. Follows NASA Spacecraft Classification.
	MAV_TYPE_SPACECRAFT_ORBITER MAV_TYPE = minimal.MAV_TYPE_SPACECRAFT_ORBITER
	// A generic four-legged ground vehicle (e.g., a robot dog).
	MAV_TYPE_GROUND_QUADRUPED MAV_TYPE = minimal.MAV_TYPE_GROUND_QUADRUPED
	// VTOL hybrid of helicopter and autogyro. It has a main rotor for lift and separate propellers for forward flight. The rotor must be powered for hover but can autorotate in cruise flight. See: https://en.wikipedia.org/wiki/Gyrodyne
	MAV_TYPE_VTOL_GYRODYNE MAV_TYPE = minimal.MAV_TYPE_VTOL_GYRODYNE
	// Gripper
	MAV_TYPE_GRIPPER MAV_TYPE = minimal.MAV_TYPE_GRIPPER
)
