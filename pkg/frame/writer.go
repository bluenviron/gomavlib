package frame

import (
	"fmt"
	"io"
	"time"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
	"github.com/bluenviron/gomavlib/v4/pkg/message"
)

func encodeMessageInFrame(fr Frame, mp *message.ReadWriter) {
	_, isV2 := fr.(*V2Frame)
	msgRaw := mp.Write(fr.GetMessage(), isV2)

	switch ff := fr.(type) {
	case *V1Frame:
		ff.Message = msgRaw
	case *V2Frame:
		ff.Message = msgRaw
	}
}

// Writer is a Frame writer.
type Writer struct {
	// underlying byte writer.
	ByteWriter io.Writer

	// (optional) dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	//
	// private
	//

	timeNow func() time.Time
	bw      []byte
}

// Initialize allocates a Writer.
func (w *Writer) Initialize() error {
	if w.ByteWriter == nil {
		return fmt.Errorf("ByteWriter not provided")
	}

	if w.timeNow == nil {
		w.timeNow = time.Now
	}

	w.bw = make([]byte, bufferSize)

	return nil
}

// Write writes a Frame.
// It must not be called by multiple routines in parallel.
func (w *Writer) Write(fr Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if w.DialectRW == nil {
			return fmt.Errorf("dialect is nil")
		}

		mp := w.DialectRW.GetMessage(fr.GetMessage().GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		encodeMessageInFrame(fr, mp)
	}

	return w.writeFrameInner(fr)
}

func (w *Writer) writeFrameInner(fr Frame) error {
	n, err := fr.marshalTo(w.bw, fr.GetMessage().(*message.MessageRaw).Payload)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.ByteWriter.Write(w.bw[:n])
	return err
}
