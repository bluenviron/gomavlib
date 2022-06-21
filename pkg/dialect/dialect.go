// Package dialect contains the dialect definition and parser.
package dialect

import (
	"github.com/aler9/gomavlib/pkg/msg"
)

// Dialect is a Mavlink dialect.
type Dialect struct {
	// Version is the dialect version.
	Version int

	// Messages contains the messages of the dialect.
	Messages []msg.Message
}
