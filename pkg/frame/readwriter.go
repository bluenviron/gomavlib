package frame

import (
	"bufio"
	"fmt"
	"io"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
)

// ReadWriter is a Frame Reader and Writer.
type ReadWriter struct {
	// underlying byte ReadWriter.
	ByteReadWriter io.ReadWriter

	// (optional) dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key

	*Reader
	*Writer
}

// Initialize initializes ReadWriter.
func (rw *ReadWriter) Initialize() error {
	if rw.ByteReadWriter == nil {
		return fmt.Errorf("ByteReadWriter not provided")
	}

	r := &Reader{
		BufByteReader: bufio.NewReaderSize(rw.ByteReadWriter, bufferSize),
		DialectRW:     rw.DialectRW,
		InKey:         rw.InKey,
	}
	err := r.Initialize()
	if err != nil {
		return err
	}

	w := &Writer{
		ByteWriter: rw.ByteReadWriter,
		DialectRW:  rw.DialectRW,
	}
	err = w.Initialize()
	if err != nil {
		return err
	}

	rw.Reader = r
	rw.Writer = w

	return nil
}
