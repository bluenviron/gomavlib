// Package tlog contains a Telemetry log reader and writer.
package tlog

import (
	"bufio"
	"io"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
)

// Entry is a Telemetry log entry.
type Entry struct {
	// Timestamp of the entry.
	Time time.Time

	// Mavlink frame.
	Frame frame.Frame
}

// Reader is a telemetry log reader.
// Specification: https://docs.qgroundcontrol.com/master/en/qgc-dev-guide/file_formats/mavlink.html
type Reader struct {
	// underlying byte reader.
	ByteReader io.Reader

	// (optional) dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	//
	// private
	//

	br          *bufio.Reader
	frameReader *frame.Reader
}

// Initialize initializes Reader.
func (r *Reader) Initialize() error {
	r.br = bufio.NewReader(r.ByteReader)

	r.frameReader = &frame.Reader{
		BufByteReader: r.br,
		DialectRW:     r.DialectRW,
	}
	err := r.frameReader.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (r *Reader) Read() (*Entry, error) {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r.br, buf)
	if err != nil {
		return nil, err
	}

	epoch := int64(buf[0])<<56 |
		int64(buf[1])<<48 |
		int64(buf[2])<<40 |
		int64(buf[3])<<32 |
		int64(buf[4])<<24 |
		int64(buf[5])<<16 |
		int64(buf[6])<<8 |
		int64(buf[7])
	t := time.Unix(epoch/1000000, (epoch%1000000)*1000).UTC()

	fr, err := r.frameReader.Read()
	if err != nil {
		return nil, err
	}

	return &Entry{
		Time:  t,
		Frame: fr,
	}, nil
}
