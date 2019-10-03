// Autogenerated with dialgen, do not edit.
//
// Generated from revision https://github.com/mavlink/mavlink/tree/1a1534ef64a75dbbdceaaf462c7d963b2be5a321
//
package minimal

import (
	"github.com/gswly/gomavlib"
)

// Dialect contains the dialect object that can be passed to the library.
var Dialect = dialect

// dialect is not exposed directly such that it is not displayed in godoc.
var dialect = gomavlib.MustDialect(3, []gomavlib.Message{
	// minimal.xml
	&MessageHeartbeat{},
	&MessageProtocolVersion{},
})

// Micro air vehicle / autopilot classes. This identifies the individual model.
type MAV_AUTOPILOT int

const (
	// Generic autopilot, full support for everything
	MAV_AUTOPILOT_GENERIC MAV_AUTOPILOT = 0
	// Reserved for future use.
	MAV_AUTOPILOT_RESERVED MAV_AUTOPILOT = 1
	// SLUGS autopilot, http://slugsuav.soe.ucsc.edu
	MAV_AUTOPILOT_SLUGS MAV_AUTOPILOT = 2
	// ArduPilot - Plane/Copter/Rover/Sub/Tracker, http://ardupilot.org
	MAV_AUTOPILOT_ARDUPILOTMEGA MAV_AUTOPILOT = 3
	// OpenPilot, http://openpilot.org
	MAV_AUTOPILOT_OPENPILOT MAV_AUTOPILOT = 4
	// Generic autopilot only supporting simple waypoints
	MAV_AUTOPILOT_GENERIC_WAYPOINTS_ONLY MAV_AUTOPILOT = 5
	// Generic autopilot supporting waypoints and other simple navigation commands
	MAV_AUTOPILOT_GENERIC_WAYPOINTS_AND_SIMPLE_NAVIGATION_ONLY MAV_AUTOPILOT = 6
	// Generic autopilot supporting the full mission command set
	MAV_AUTOPILOT_GENERIC_MISSION_FULL MAV_AUTOPILOT = 7
	// No valid autopilot, e.g. a GCS or other MAVLink component
	MAV_AUTOPILOT_INVALID MAV_AUTOPILOT = 8
	// PPZ UAV - http://nongnu.org/paparazzi
	MAV_AUTOPILOT_PPZ MAV_AUTOPILOT = 9
	// UAV Dev Board
	MAV_AUTOPILOT_UDB MAV_AUTOPILOT = 10
	// FlexiPilot
	MAV_AUTOPILOT_FP MAV_AUTOPILOT = 11
	// PX4 Autopilot - http://px4.io/
	MAV_AUTOPILOT_PX4 MAV_AUTOPILOT = 12
	// SMACCMPilot - http://smaccmpilot.org
	MAV_AUTOPILOT_SMACCMPILOT MAV_AUTOPILOT = 13
	// AutoQuad -- http://autoquad.org
	MAV_AUTOPILOT_AUTOQUAD MAV_AUTOPILOT = 14
	// Armazila -- http://armazila.com
	MAV_AUTOPILOT_ARMAZILA MAV_AUTOPILOT = 15
	// Aerob -- http://aerob.ru
	MAV_AUTOPILOT_AEROB MAV_AUTOPILOT = 16
	// ASLUAV autopilot -- http://www.asl.ethz.ch
	MAV_AUTOPILOT_ASLUAV MAV_AUTOPILOT = 17
	// SmartAP Autopilot - http://sky-drones.com
	MAV_AUTOPILOT_SMARTAP MAV_AUTOPILOT = 18
	// AirRails - http://uaventure.com
	MAV_AUTOPILOT_AIRRAILS MAV_AUTOPILOT = 19
)

// Commands to be executed by the MAV. They can be executed on user request, or as part of a mission script. If the action is used in a mission, the parameter mapping to the waypoint/mission message is as follows: Param 1, Param 2, Param 3, Param 4, X: Param 5, Y:Param 6, Z:Param 7. This command list is similar what ARINC 424 is for commercial aircraft: A data format how to interpret waypoint/mission data. See https://mavlink.io/en/guide/xml_schema.html#MAV_CMD for information about the structure of the MAV_CMD entries
type MAV_CMD int

const (
	// Request MAVLink protocol version compatibility
	MAV_CMD_REQUEST_PROTOCOL_VERSION MAV_CMD = 519
)

// Component ids (values) for the different types and instances of onboard hardware/software that might make up a MAVLink system (autopilot, cameras, servos, GPS systems, avoidance systems etc.).      Components must use the appropriate ID in their source address when sending messages. Components can also use IDs to determine if they are the intended recipient of an incoming message. The MAV_COMP_ID_ALL value is used to indicate messages that must be processed by all components.      When creating new entries, components that can have multiple instances (e.g. cameras, servos etc.) should be allocated sequential values. An appropriate number of values should be left free after these components to allow the number of instances to be expanded.
type MAV_COMPONENT int

const (
	// Used to broadcast messages to all components of the receiving system. Components should attempt to process messages with this component ID and forward to components on any other interfaces.
	MAV_COMP_ID_ALL MAV_COMPONENT = 0
	// System flight controller component ("autopilot"). Only one autopilot is expected in a particular system.
	MAV_COMP_ID_AUTOPILOT1 MAV_COMPONENT = 1
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER1 MAV_COMPONENT = 25
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER2 MAV_COMPONENT = 26
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER3 MAV_COMPONENT = 27
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER4 MAV_COMPONENT = 28
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER5 MAV_COMPONENT = 29
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER6 MAV_COMPONENT = 30
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER7 MAV_COMPONENT = 31
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER8 MAV_COMPONENT = 32
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER9 MAV_COMPONENT = 33
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER10 MAV_COMPONENT = 34
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER11 MAV_COMPONENT = 35
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER12 MAV_COMPONENT = 36
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER13 MAV_COMPONENT = 37
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER14 MAV_COMPONENT = 38
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER15 MAV_COMPONENT = 39
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USE16 MAV_COMPONENT = 40
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER17 MAV_COMPONENT = 41
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER18 MAV_COMPONENT = 42
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER19 MAV_COMPONENT = 43
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER20 MAV_COMPONENT = 44
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER21 MAV_COMPONENT = 45
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER22 MAV_COMPONENT = 46
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER23 MAV_COMPONENT = 47
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER24 MAV_COMPONENT = 48
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER25 MAV_COMPONENT = 49
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER26 MAV_COMPONENT = 50
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER27 MAV_COMPONENT = 51
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER28 MAV_COMPONENT = 52
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER29 MAV_COMPONENT = 53
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER30 MAV_COMPONENT = 54
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER31 MAV_COMPONENT = 55
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER32 MAV_COMPONENT = 56
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER33 MAV_COMPONENT = 57
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER34 MAV_COMPONENT = 58
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER35 MAV_COMPONENT = 59
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER36 MAV_COMPONENT = 60
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER37 MAV_COMPONENT = 61
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER38 MAV_COMPONENT = 62
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER39 MAV_COMPONENT = 63
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER40 MAV_COMPONENT = 64
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER41 MAV_COMPONENT = 65
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER42 MAV_COMPONENT = 66
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER43 MAV_COMPONENT = 67
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER44 MAV_COMPONENT = 68
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER45 MAV_COMPONENT = 69
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER46 MAV_COMPONENT = 70
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER47 MAV_COMPONENT = 71
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER48 MAV_COMPONENT = 72
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER49 MAV_COMPONENT = 73
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER50 MAV_COMPONENT = 74
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER51 MAV_COMPONENT = 75
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER52 MAV_COMPONENT = 76
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER53 MAV_COMPONENT = 77
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER54 MAV_COMPONENT = 78
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER55 MAV_COMPONENT = 79
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER56 MAV_COMPONENT = 80
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER57 MAV_COMPONENT = 81
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER58 MAV_COMPONENT = 82
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER59 MAV_COMPONENT = 83
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER60 MAV_COMPONENT = 84
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER61 MAV_COMPONENT = 85
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER62 MAV_COMPONENT = 86
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER63 MAV_COMPONENT = 87
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER64 MAV_COMPONENT = 88
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER65 MAV_COMPONENT = 89
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER66 MAV_COMPONENT = 90
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER67 MAV_COMPONENT = 91
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER68 MAV_COMPONENT = 92
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER69 MAV_COMPONENT = 93
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER70 MAV_COMPONENT = 94
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER71 MAV_COMPONENT = 95
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER72 MAV_COMPONENT = 96
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER73 MAV_COMPONENT = 97
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER74 MAV_COMPONENT = 98
	// Id for a component on privately managed MAVLink network. Can be used for any purpose but may not be published by components outside of the private network.
	MAV_COMP_ID_USER75 MAV_COMPONENT = 99
	// Camera #1.
	MAV_COMP_ID_CAMERA MAV_COMPONENT = 100
	// Camera #2.
	MAV_COMP_ID_CAMERA2 MAV_COMPONENT = 101
	// Camera #3.
	MAV_COMP_ID_CAMERA3 MAV_COMPONENT = 102
	// Camera #4.
	MAV_COMP_ID_CAMERA4 MAV_COMPONENT = 103
	// Camera #5.
	MAV_COMP_ID_CAMERA5 MAV_COMPONENT = 104
	// Camera #6.
	MAV_COMP_ID_CAMERA6 MAV_COMPONENT = 105
	// Servo #1.
	MAV_COMP_ID_SERVO1 MAV_COMPONENT = 140
	// Servo #2.
	MAV_COMP_ID_SERVO2 MAV_COMPONENT = 141
	// Servo #3.
	MAV_COMP_ID_SERVO3 MAV_COMPONENT = 142
	// Servo #4.
	MAV_COMP_ID_SERVO4 MAV_COMPONENT = 143
	// Servo #5.
	MAV_COMP_ID_SERVO5 MAV_COMPONENT = 144
	// Servo #6.
	MAV_COMP_ID_SERVO6 MAV_COMPONENT = 145
	// Servo #7.
	MAV_COMP_ID_SERVO7 MAV_COMPONENT = 146
	// Servo #8.
	MAV_COMP_ID_SERVO8 MAV_COMPONENT = 147
	// Servo #9.
	MAV_COMP_ID_SERVO9 MAV_COMPONENT = 148
	// Servo #10.
	MAV_COMP_ID_SERVO10 MAV_COMPONENT = 149
	// Servo #11.
	MAV_COMP_ID_SERVO11 MAV_COMPONENT = 150
	// Servo #12.
	MAV_COMP_ID_SERVO12 MAV_COMPONENT = 151
	// Servo #13.
	MAV_COMP_ID_SERVO13 MAV_COMPONENT = 152
	// Servo #14.
	MAV_COMP_ID_SERVO14 MAV_COMPONENT = 153
	// Gimbal component.
	MAV_COMP_ID_GIMBAL MAV_COMPONENT = 154
	// Logging component.
	MAV_COMP_ID_LOG MAV_COMPONENT = 155
	// Automatic Dependent Surveillance-Broadcast (ADS-B) component.
	MAV_COMP_ID_ADSB MAV_COMPONENT = 156
	// On Screen Display (OSD) devices for video links.
	MAV_COMP_ID_OSD MAV_COMPONENT = 157
	// Generic autopilot peripheral component ID. Meant for devices that do not implement the parameter microservice.
	MAV_COMP_ID_PERIPHERAL MAV_COMPONENT = 158
	// Gimbal ID for QX1.
	MAV_COMP_ID_QX1_GIMBAL MAV_COMPONENT = 159
	// FLARM collision alert component.
	MAV_COMP_ID_FLARM MAV_COMPONENT = 160
	// Component that can generate/supply a mission flight plan (e.g. GCS or developer API).
	MAV_COMP_ID_MISSIONPLANNER MAV_COMPONENT = 190
	// Component that finds an optimal path between points based on a certain constraint (e.g. minimum snap, shortest path, cost, etc.).
	MAV_COMP_ID_PATHPLANNER MAV_COMPONENT = 195
	// Component that plans a collision free path between two points.
	MAV_COMP_ID_OBSTACLE_AVOIDANCE MAV_COMPONENT = 196
	// Component that provides position estimates using VIO techniques.
	MAV_COMP_ID_VISUAL_INERTIAL_ODOMETRY MAV_COMPONENT = 197
	// Inertial Measurement Unit (IMU) #1.
	MAV_COMP_ID_IMU MAV_COMPONENT = 200
	// Inertial Measurement Unit (IMU) #2.
	MAV_COMP_ID_IMU_2 MAV_COMPONENT = 201
	// Inertial Measurement Unit (IMU) #3.
	MAV_COMP_ID_IMU_3 MAV_COMPONENT = 202
	// GPS #1.
	MAV_COMP_ID_GPS MAV_COMPONENT = 220
	// GPS #2.
	MAV_COMP_ID_GPS2 MAV_COMPONENT = 221
	// Component to bridge MAVLink to UDP (i.e. from a UART).
	MAV_COMP_ID_UDP_BRIDGE MAV_COMPONENT = 240
	// Component to bridge to UART (i.e. from UDP).
	MAV_COMP_ID_UART_BRIDGE MAV_COMPONENT = 241
	// Component for handling system messages (e.g. to ARM, takeoff, etc.).
	MAV_COMP_ID_SYSTEM_CONTROL MAV_COMPONENT = 250
)

// These flags encode the MAV mode.
type MAV_MODE_FLAG int

const (
	// 0b10000000 MAV safety set to armed. Motors are enabled / running / can start. Ready to fly. Additional note: this flag is to be ignore when sent in the command MAV_CMD_DO_SET_MODE and MAV_CMD_COMPONENT_ARM_DISARM shall be used instead. The flag can still be used to report the armed state.
	MAV_MODE_FLAG_SAFETY_ARMED MAV_MODE_FLAG = 128
	// 0b01000000 remote control input is enabled.
	MAV_MODE_FLAG_MANUAL_INPUT_ENABLED MAV_MODE_FLAG = 64
	// 0b00100000 hardware in the loop simulation. All motors / actuators are blocked, but internal software is full operational.
	MAV_MODE_FLAG_HIL_ENABLED MAV_MODE_FLAG = 32
	// 0b00010000 system stabilizes electronically its attitude (and optionally position). It needs however further control inputs to move around.
	MAV_MODE_FLAG_STABILIZE_ENABLED MAV_MODE_FLAG = 16
	// 0b00001000 guided mode enabled, system flies waypoints / mission items.
	MAV_MODE_FLAG_GUIDED_ENABLED MAV_MODE_FLAG = 8
	// 0b00000100 autonomous mode enabled, system finds its own goal positions. Guided flag can be set or not, depends on the actual implementation.
	MAV_MODE_FLAG_AUTO_ENABLED MAV_MODE_FLAG = 4
	// 0b00000010 system has a test mode enabled. This flag is intended for temporary system tests and should not be used for stable implementations.
	MAV_MODE_FLAG_TEST_ENABLED MAV_MODE_FLAG = 2
	// 0b00000001 Reserved for future use.
	MAV_MODE_FLAG_CUSTOM_MODE_ENABLED MAV_MODE_FLAG = 1
)

// These values encode the bit positions of the decode position. These values can be used to read the value of a flag bit by combining the base_mode variable with AND with the flag position value. The result will be either 0 or 1, depending on if the flag is set or not.
type MAV_MODE_FLAG_DECODE_POSITION int

const (
	// First bit:  10000000
	MAV_MODE_FLAG_DECODE_POSITION_SAFETY MAV_MODE_FLAG_DECODE_POSITION = 128
	// Second bit: 01000000
	MAV_MODE_FLAG_DECODE_POSITION_MANUAL MAV_MODE_FLAG_DECODE_POSITION = 64
	// Third bit:  00100000
	MAV_MODE_FLAG_DECODE_POSITION_HIL MAV_MODE_FLAG_DECODE_POSITION = 32
	// Fourth bit: 00010000
	MAV_MODE_FLAG_DECODE_POSITION_STABILIZE MAV_MODE_FLAG_DECODE_POSITION = 16
	// Fifth bit:  00001000
	MAV_MODE_FLAG_DECODE_POSITION_GUIDED MAV_MODE_FLAG_DECODE_POSITION = 8
	// Sixth bit:   00000100
	MAV_MODE_FLAG_DECODE_POSITION_AUTO MAV_MODE_FLAG_DECODE_POSITION = 4
	// Seventh bit: 00000010
	MAV_MODE_FLAG_DECODE_POSITION_TEST MAV_MODE_FLAG_DECODE_POSITION = 2
	// Eighth bit: 00000001
	MAV_MODE_FLAG_DECODE_POSITION_CUSTOM_MODE MAV_MODE_FLAG_DECODE_POSITION = 1
)

//
type MAV_STATE int

const (
	// Uninitialized system, state is unknown.
	MAV_STATE_UNINIT MAV_STATE = 0
	// System is booting up.
	MAV_STATE_BOOT MAV_STATE = 1
	// System is calibrating and not flight-ready.
	MAV_STATE_CALIBRATING MAV_STATE = 2
	// System is grounded and on standby. It can be launched any time.
	MAV_STATE_STANDBY MAV_STATE = 3
	// System is active and might be already airborne. Motors are engaged.
	MAV_STATE_ACTIVE MAV_STATE = 4
	// System is in a non-normal flight mode. It can however still navigate.
	MAV_STATE_CRITICAL MAV_STATE = 5
	// System is in a non-normal flight mode. It lost control over parts or over the whole airframe. It is in mayday and going down.
	MAV_STATE_EMERGENCY MAV_STATE = 6
	// System just initialized its power-down sequence, will shut down now.
	MAV_STATE_POWEROFF MAV_STATE = 7
	// System is terminating itself.
	MAV_STATE_FLIGHT_TERMINATION MAV_STATE = 8
)

// MAVLINK component type reported in HEARTBEAT message. Flight controllers must report the type of the vehicle on which they are mounted (e.g. MAV_TYPE_OCTOROTOR). All other components must report a value appropriate for their type (e.g. a camera must use MAV_TYPE_CAMERA).
type MAV_TYPE int

const (
	// Generic micro air vehicle
	MAV_TYPE_GENERIC MAV_TYPE = 0
	// Fixed wing aircraft.
	MAV_TYPE_FIXED_WING MAV_TYPE = 1
	// Quadrotor
	MAV_TYPE_QUADROTOR MAV_TYPE = 2
	// Coaxial helicopter
	MAV_TYPE_COAXIAL MAV_TYPE = 3
	// Normal helicopter with tail rotor.
	MAV_TYPE_HELICOPTER MAV_TYPE = 4
	// Ground installation
	MAV_TYPE_ANTENNA_TRACKER MAV_TYPE = 5
	// Operator control unit / ground control station
	MAV_TYPE_GCS MAV_TYPE = 6
	// Airship, controlled
	MAV_TYPE_AIRSHIP MAV_TYPE = 7
	// Free balloon, uncontrolled
	MAV_TYPE_FREE_BALLOON MAV_TYPE = 8
	// Rocket
	MAV_TYPE_ROCKET MAV_TYPE = 9
	// Ground rover
	MAV_TYPE_GROUND_ROVER MAV_TYPE = 10
	// Surface vessel, boat, ship
	MAV_TYPE_SURFACE_BOAT MAV_TYPE = 11
	// Submarine
	MAV_TYPE_SUBMARINE MAV_TYPE = 12
	// Hexarotor
	MAV_TYPE_HEXAROTOR MAV_TYPE = 13
	// Octorotor
	MAV_TYPE_OCTOROTOR MAV_TYPE = 14
	// Tricopter
	MAV_TYPE_TRICOPTER MAV_TYPE = 15
	// Flapping wing
	MAV_TYPE_FLAPPING_WING MAV_TYPE = 16
	// Kite
	MAV_TYPE_KITE MAV_TYPE = 17
	// Onboard companion controller
	MAV_TYPE_ONBOARD_CONTROLLER MAV_TYPE = 18
	// Two-rotor VTOL using control surfaces in vertical operation in addition. Tailsitter.
	MAV_TYPE_VTOL_DUOROTOR MAV_TYPE = 19
	// Quad-rotor VTOL using a V-shaped quad config in vertical operation. Tailsitter.
	MAV_TYPE_VTOL_QUADROTOR MAV_TYPE = 20
	// Tiltrotor VTOL
	MAV_TYPE_VTOL_TILTROTOR MAV_TYPE = 21
	// VTOL reserved 2
	MAV_TYPE_VTOL_RESERVED2 MAV_TYPE = 22
	// VTOL reserved 3
	MAV_TYPE_VTOL_RESERVED3 MAV_TYPE = 23
	// VTOL reserved 4
	MAV_TYPE_VTOL_RESERVED4 MAV_TYPE = 24
	// VTOL reserved 5
	MAV_TYPE_VTOL_RESERVED5 MAV_TYPE = 25
	// Gimbal
	MAV_TYPE_GIMBAL MAV_TYPE = 26
	// ADSB system
	MAV_TYPE_ADSB MAV_TYPE = 27
	// Steerable, nonrigid airfoil
	MAV_TYPE_PARAFOIL MAV_TYPE = 28
	// Dodecarotor
	MAV_TYPE_DODECAROTOR MAV_TYPE = 29
	// Camera
	MAV_TYPE_CAMERA MAV_TYPE = 30
	// Charging station
	MAV_TYPE_CHARGING_STATION MAV_TYPE = 31
	// FLARM collision avoidance system
	MAV_TYPE_FLARM MAV_TYPE = 32
	// Servo
	MAV_TYPE_SERVO MAV_TYPE = 33
)

// minimal.xml

// The heartbeat message shows that a system or component is present and responding. The type and autopilot fields (along with the message component id), allow the receiving system to treat further messages from this system appropriately (e.g. by laying out the user interface based on the autopilot). This microservice is documented at https://mavlink.io/en/services/heartbeat.html
type MessageHeartbeat struct {
	// Type of the system (quadrotor, helicopter, etc.). Components use the same type as their associated system.
	Type MAV_TYPE `mavenum:"uint8"`
	// Autopilot type / class.
	Autopilot MAV_AUTOPILOT `mavenum:"uint8"`
	// System mode bitmap.
	BaseMode MAV_MODE_FLAG `mavenum:"uint8"`
	// A bitfield for use for autopilot-specific flags
	CustomMode uint32
	// System status flag.
	SystemStatus MAV_STATE `mavenum:"uint8"`
	// MAVLink version, not writable by user, gets added by protocol because of magic data type: uint8_t_mavlink_version
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetId() uint32 {
	return 0
}

// Version and capability of protocol version. This message is the response to REQUEST_PROTOCOL_VERSION and is used as part of the handshaking to establish which MAVLink version should be used on the network. Every node should respond to REQUEST_PROTOCOL_VERSION to enable the handshaking. Library implementers should consider adding this into the default decoding state machine to allow the protocol core to respond directly.
type MessageProtocolVersion struct {
	// Currently active MAVLink version number * 100: v1.0 is 100, v2.0 is 200, etc.
	Version uint16
	// Minimum MAVLink version supported
	MinVersion uint16
	// Maximum MAVLink version supported (set to the same value as version by default)
	MaxVersion uint16
	// The first 8 bytes (not characters printed in hex!) of the git hash.
	SpecVersionHash [8]uint8
	// The first 8 bytes (not characters printed in hex!) of the git hash.
	LibraryVersionHash [8]uint8
}

func (*MessageProtocolVersion) GetId() uint32 {
	return 300
}
