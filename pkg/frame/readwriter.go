package frame

import (
	"io"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
)

// ReadWriterConf is the configuration of a ReadWriter.
type ReadWriterConf struct {
	// the underlying bytes ReadWriter.
	ReadWriter io.ReadWriter

	// (optional) the dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key
}

// ReadWriter is a Frame Reader and Writer.
type ReadWriter struct {
	*Reader
	*Writer
}

// NewReadWriter allocates a ReadWriter.
func NewReadWriter(conf ReadWriterConf) (*ReadWriter, error) {
	r, err := NewReader(ReaderConf{
		Reader:    conf.ReadWriter,
		DialectRW: conf.DialectRW,
		InKey:     conf.InKey,
	})
	if err != nil {
		return nil, err
	}

	w, err := NewWriter(WriterConf{
		Writer:             conf.ReadWriter,
		DialectRW:          conf.DialectRW,
		OutVersion:         conf.OutVersion,
		OutSystemID:        conf.OutSystemID,
		OutComponentID:     conf.OutComponentID,
		OutSignatureLinkID: conf.OutSignatureLinkID,
		OutKey:             conf.OutKey,
	})
	if err != nil {
		return nil, err
	}

	return &ReadWriter{
		Reader: r,
		Writer: w,
	}, nil
}
