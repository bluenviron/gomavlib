//autogenerated:yes
//nolint:revive,misspell,govet,lll
package all

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/development"
)

// RC channel outputs from a MAVLink RC receiver for input to a flight controller or other components (allows an RC receiver to connect via MAVLink instead of some other protocol such as PPM-Sum or S.BUS).
// Note that this is not intended to be an over-the-air format, and does not replace RC_CHANNELS and similar messages reported by the flight controller.
// The target_system field should normally be set to the system id of the system to control, typically the flight controller.
// The target_component field can normally be set to 0, so that all components of the system can receive the message.
// The channels array field can publish up to 32 channels; the number of channel items used in the array is specified in the count field.
// The time_last_update_ms field contains the timestamp of the last received valid channels data in the receiver's time domain.
// The count field indicates the first index of the channel array that is not used for channel data (this and later indexes are zero-filled).
// The RADIO_RC_CHANNELS_FLAGS_OUTDATED flag is set by the receiver if the channels data is not up-to-date (for example, if new data from the transmitter could not be validated so the last valid data is resent).
// The RADIO_RC_CHANNELS_FLAGS_FAILSAFE failsafe flag is set by the receiver if the receiver's failsafe condition is met (implementation dependent, e.g., connection to the RC radio is lost).
// In this case time_last_update_ms still contains the timestamp of the last valid channels data, but the content of the channels data is not defined by the protocol (it is up to the implementation of the receiver).
// For instance, the channels data could contain failsafe values configured in the receiver; the default is to carry the last valid data.
// Note: The RC channels fields are extensions to ensure that they are located at the end of the serialized payload and subject to MAVLink's trailing-zero trimming.
type MessageRadioRcChannels = development.MessageRadioRcChannels
