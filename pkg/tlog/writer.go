package tlog

import (
	"io"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

// Writer is a telemetry log writer.
// Specification: https://docs.qgroundcontrol.com/master/en/qgc-dev-guide/file_formats/mavlink.html
type Writer struct {
	// underlying byte writer.
	ByteWriter io.Writer

	// (optional) dialect which contains the messages that will be written.
	// If not provided, messages are encoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	//
	// private
	//

	frameWriter *frame.Writer
}

// Initialize initializes Writer.
func (w *Writer) Initialize() error {
	w.frameWriter = &frame.Writer{
		ByteWriter:  w.ByteWriter,
		DialectRW:   w.DialectRW,
		OutVersion:  frame.V1, // unused
		OutSystemID: 1,        // unused
	}
	err := w.frameWriter.Initialize()
	if err != nil {
		return err
	}

	return nil
}

// Write writes a telemetry log entry.
func (w *Writer) Write(entry *Entry) error {
	epoch := entry.Time.UnixMicro()
	buf := []byte{
		byte(epoch >> 56),
		byte(epoch >> 48),
		byte(epoch >> 40),
		byte(epoch >> 32),
		byte(epoch >> 24),
		byte(epoch >> 16),
		byte(epoch >> 8),
		byte(epoch),
	}
	_, err := w.ByteWriter.Write(buf)
	if err != nil {
		return err
	}

	err = w.frameWriter.WriteFrame(entry.Frame)
	if err != nil {
		return err
	}

	return nil
}
