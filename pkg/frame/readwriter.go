package frame

import (
	"io"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
)

// ReadWriterConf is the configuration of a ReadWriter.
//
// Deprecated: configuration has been moved inside ReadWriter.
type ReadWriterConf struct {
	// underlying bytes ReadWriter.
	ReadWriter io.ReadWriter

	// (optional) dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key
}

// NewReadWriter allocates a ReadWriter.
//
// Deprecated: replaced by ReadWriter.Initialize().
func NewReadWriter(conf ReadWriterConf) (*ReadWriter, error) {
	rw := &ReadWriter{
		ByteReadWriter:     conf.ReadWriter,
		DialectRW:          conf.DialectRW,
		InKey:              conf.InKey,
		OutVersion:         conf.OutVersion,
		OutSystemID:        conf.OutSystemID,
		OutComponentID:     conf.OutComponentID,
		OutSignatureLinkID: conf.OutSignatureLinkID,
		OutKey:             conf.OutKey,
	}
	err := rw.Initialize()
	return rw, err
}

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

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key

	*Reader
	*Writer
}

// Initialize initializes ReadWriter.
func (rw *ReadWriter) Initialize() error {
	r, err := NewReader(ReaderConf{
		Reader:    rw.ByteReadWriter,
		DialectRW: rw.DialectRW,
		InKey:     rw.InKey,
	})
	if err != nil {
		return err
	}

	w, err := NewWriter(WriterConf{
		Writer:             rw.ByteReadWriter,
		DialectRW:          rw.DialectRW,
		OutVersion:         rw.OutVersion,
		OutSystemID:        rw.OutSystemID,
		OutComponentID:     rw.OutComponentID,
		OutSignatureLinkID: rw.OutSignatureLinkID,
		OutKey:             rw.OutKey,
	})
	if err != nil {
		return err
	}

	rw.Reader = r
	rw.Writer = w

	return nil
}
