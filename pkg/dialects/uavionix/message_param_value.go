//autogenerated:yes
//nolint:revive,misspell,govet,lll
package uavionix

import (
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
)

// Emit the value of a onboard parameter. The inclusion of param_count and param_index in the message allows the recipient to keep track of received parameters and allows him to re-request missing parameters after a loss or timeout. The parameter microservice is documented at https://mavlink.io/en/services/parameter.html
type MessageParamValue = common.MessageParamValue
