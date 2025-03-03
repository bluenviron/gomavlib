package dialect

import (
	"fmt"

	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// NewReadWriter allocates a ReadWriter.
//
// Deprecated: replaced by ReadWriter.Initialize().
func NewReadWriter(d *Dialect) (*ReadWriter, error) {
	rw := &ReadWriter{Dialect: d}
	err := rw.Initialize()
	return rw, err
}

// ReadWriter is a Dialect Reader and Writer.
type ReadWriter struct {
	Dialect *Dialect

	messageRWs map[uint32]*message.ReadWriter
}

// Initialize initializes a ReadWriter.
func (rw *ReadWriter) Initialize() error {
	rw.messageRWs = make(map[uint32]*message.ReadWriter)

	for _, m := range rw.Dialect.Messages {
		if _, ok := rw.messageRWs[m.GetID()]; ok {
			return fmt.Errorf("duplicate message with id %d", m.GetID())
		}

		de := &message.ReadWriter{Message: m}
		err := de.Initialize()
		if err != nil {
			return fmt.Errorf("message %T: %w", m, err)
		}

		rw.messageRWs[m.GetID()] = de
	}

	return nil
}

// GetMessage returns the ReadWriter of a message.
func (rw *ReadWriter) GetMessage(id uint32) *message.ReadWriter {
	mrw, ok := rw.messageRWs[id]
	if !ok {
		return nil
	}
	return mrw
}
