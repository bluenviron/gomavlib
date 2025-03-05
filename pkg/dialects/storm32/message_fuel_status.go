//autogenerated:yes
//nolint:revive,misspell,govet,lll
package storm32

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Fuel status.
// This message provides "generic" fuel level information for  in a GCS and for triggering failsafes in an autopilot.
// The fuel type and associated units for fields in this message are defined in the enum MAV_FUEL_TYPE.
// The reported `consumed_fuel` and `remaining_fuel` must only be supplied if measured: they must not be inferred from the `maximum_fuel` and the other value.
// A recipient can assume that if these fields are supplied they are accurate.
// If not provided, the recipient can infer `remaining_fuel` from `maximum_fuel` and `consumed_fuel` on the assumption that the fuel was initially at its maximum (this is what battery monitors assume).
// Note however that this is an assumption, and the UI should prompt the user appropriately (i.e. notify user that they should fill the tank before boot).
// This kind of information may also be sent in fuel-specific messages such as BATTERY_STATUS_V2.
// If both messages are sent for the same fuel system, the ids and corresponding information must match.
// This should be streamed (nominally at 0.1 Hz).
type MessageFuelStatus = common.MessageFuelStatus
