package dialect

import (
	"fmt"

	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

// ReadWriter is a Dialect Reader and Writer.
type ReadWriter struct {
	messageRWs map[uint32]*message.ReadWriter
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(d *Dialect) (*ReadWriter, error) {
	rw := &ReadWriter{
		messageRWs: make(map[uint32]*message.ReadWriter),
	}

	for _, m := range d.Messages {
		if _, ok := rw.messageRWs[m.GetID()]; ok {
			return nil, fmt.Errorf("duplicate message with id %d", m.GetID())
		}

		de, err := message.NewReadWriter(m)
		if err != nil {
			return nil, fmt.Errorf("message %T: %w", m, err)
		}

		rw.messageRWs[m.GetID()] = de
	}

	return rw, nil
}

// GetMessage returns the ReadWriter of a message.
func (rw *ReadWriter) GetMessage(id uint32) *message.ReadWriter {
	mrw, ok := rw.messageRWs[id]
	if !ok {
		return nil
	}
	return mrw
}
