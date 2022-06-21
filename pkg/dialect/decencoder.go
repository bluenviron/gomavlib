package dialect

import (
	"fmt"

	"github.com/aler9/gomavlib/pkg/message"
)

// DecEncoder is an object that allows to decode and encode a Dialect.
type DecEncoder struct {
	MessageDEs map[uint32]*message.DecEncoder
}

// NewDecEncoder allocates a DecEncoder.
func NewDecEncoder(d *Dialect) (*DecEncoder, error) {
	dde := &DecEncoder{
		MessageDEs: make(map[uint32]*message.DecEncoder),
	}

	for _, m := range d.Messages {
		if _, ok := dde.MessageDEs[m.GetID()]; ok {
			return nil, fmt.Errorf("duplicate message with id %d", m.GetID())
		}

		de, err := message.NewDecEncoder(m)
		if err != nil {
			return nil, fmt.Errorf("message %T: %s", m, err)
		}

		dde.MessageDEs[m.GetID()] = de
	}

	return dde, nil
}
