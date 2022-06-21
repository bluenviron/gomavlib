package dialect

import (
	"fmt"

	"github.com/aler9/gomavlib/pkg/message"
)

// ReadWriter is a Dialect Reader and Writer.
type ReadWriter struct {
	MessageDEs map[uint32]*message.ReadWriter
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(d *Dialect) (*ReadWriter, error) {
	dde := &ReadWriter{
		MessageDEs: make(map[uint32]*message.ReadWriter),
	}

	for _, m := range d.Messages {
		if _, ok := dde.MessageDEs[m.GetID()]; ok {
			return nil, fmt.Errorf("duplicate message with id %d", m.GetID())
		}

		de, err := message.NewReadWriter(m)
		if err != nil {
			return nil, fmt.Errorf("message %T: %s", m, err)
		}

		dde.MessageDEs[m.GetID()] = de
	}

	return dde, nil
}
