package gomavlib

import (
	"fmt"

	"github.com/aler9/gomavlib/msg"
)

// Dialect is a dialect; it contains available messages and the
// configuration needed to encode and decode them.
type Dialect struct {
	version  uint
	messageDEs map[uint32]*msg.DecEncoder
}

// NewDialect allocates a Dialect.
func NewDialect(version uint, messages []msg.Message) (*Dialect, error) {
	d := &Dialect{
		version:  version,
		messageDEs: make(map[uint32]*msg.DecEncoder),
	}

	for _, m := range messages {
		de, err := msg.NewDecEncoder(m)
		if err != nil {
			return nil, fmt.Errorf("message %T: %s", m, err)
		}
		d.messageDEs[m.GetId()] = de
	}

	return d, nil
}

// MustDialect is like NewDialect but panics in case of error.
func MustDialect(version uint, messages []msg.Message) *Dialect {
	d, err := NewDialect(version, messages)
	if err != nil {
		panic(err)
	}
	return d
}
