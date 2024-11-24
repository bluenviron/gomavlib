// Package dialect contains the dialect definition and its parser.
package dialect

import (
	"github.com/chrisdalke/gomavlib/v3/pkg/message"
)

// Dialect is a Mavlink dialect.
type Dialect struct {
	// Version is the dialect version.
	Version int

	// Messages contains the messages of the dialect.
	Messages []message.Message
}
