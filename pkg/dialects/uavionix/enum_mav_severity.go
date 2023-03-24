//autogenerated:yes
//nolint:revive,misspell,govet,lll,dupl,gocritic
package uavionix

import (
	"github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

// Indicates the severity level, generally used for status messages to indicate their relative urgency. Based on RFC-5424 using expanded definitions at: http://www.kiwisyslog.com/kb/info:-syslog-message-levels/.
type MAV_SEVERITY = common.MAV_SEVERITY

const (
	// System is unusable. This is a "panic" condition.
	MAV_SEVERITY_EMERGENCY MAV_SEVERITY = common.MAV_SEVERITY_EMERGENCY
	// Action should be taken immediately. Indicates error in non-critical systems.
	MAV_SEVERITY_ALERT MAV_SEVERITY = common.MAV_SEVERITY_ALERT
	// Action must be taken immediately. Indicates failure in a primary system.
	MAV_SEVERITY_CRITICAL MAV_SEVERITY = common.MAV_SEVERITY_CRITICAL
	// Indicates an error in secondary/redundant systems.
	MAV_SEVERITY_ERROR MAV_SEVERITY = common.MAV_SEVERITY_ERROR
	// Indicates about a possible future error if this is not resolved within a given timeframe. Example would be a low battery warning.
	MAV_SEVERITY_WARNING MAV_SEVERITY = common.MAV_SEVERITY_WARNING
	// An unusual event has occurred, though not an error condition. This should be investigated for the root cause.
	MAV_SEVERITY_NOTICE MAV_SEVERITY = common.MAV_SEVERITY_NOTICE
	// Normal operational messages. Useful for logging. No action is required for these messages.
	MAV_SEVERITY_INFO MAV_SEVERITY = common.MAV_SEVERITY_INFO
	// Useful non-operational messages that can assist in debugging. These should not occur during normal operation.
	MAV_SEVERITY_DEBUG MAV_SEVERITY = common.MAV_SEVERITY_DEBUG
)
