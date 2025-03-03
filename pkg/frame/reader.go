package frame

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

const (
	bufferSize = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 0o1, 0o1, 0, 0, 0, 0, time.UTC)

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) ReadError {
	return ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

// ReaderConf is the configuration of a Reader.
//
// Deprecated: configuration has been moved into Reader.
type ReaderConf struct {
	// underlying bytes reader.
	Reader io.Reader

	// (optional) dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key
}

// NewReader allocates a Reader.
//
// Deprecated: replaced by Reader.Initialize().
func NewReader(conf ReaderConf) (*Reader, error) {
	r := &Reader{
		ByteReader: conf.Reader,
		DialectRW:  conf.DialectRW,
		InKey:      conf.InKey,
	}
	err := r.Initialize()
	return r, err
}

// Reader is a Frame reader.
type Reader struct {
	// underlying byte reader.
	BufByteReader *bufio.Reader

	// underlying byte reader.
	//
	// Deprecated: replaced by BufByteReader
	ByteReader io.Reader

	// (optional) dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectRW *dialect.ReadWriter

	// (optional) secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key

	//
	// private
	//

	curReadSignatureTime uint64
}

// Initialize initializes a Reader.
func (r *Reader) Initialize() error {
	if r.ByteReader != nil {
		r.BufByteReader = bufio.NewReaderSize(r.ByteReader, bufferSize)
	}

	if r.BufByteReader == nil {
		return fmt.Errorf("BufByteReader not provided")
	}

	return nil
}

// Read reads a Frame from the reader.
// It must not be called by multiple routines in parallel.
func (r *Reader) Read() (Frame, error) {
	magicByte, err := r.BufByteReader.ReadByte()
	if err != nil {
		return nil, err
	}

	var f Frame

	switch magicByte {
	case V1MagicByte:
		f = &V1Frame{}

	case V2MagicByte:
		f = &V2Frame{}

	default:
		return nil, newError("invalid magic byte: %x", magicByte)
	}

	err = f.unmarshal(r.BufByteReader)
	if err != nil {
		return nil, newError("%s", err.Error())
	}

	if r.InKey != nil {
		ff, ok := f.(*V2Frame)
		if !ok {
			return nil, newError("signature required but packet is not v2")
		}

		if ff.Signature == nil {
			return nil, newError("signature not present")
		}

		if sig := ff.GenerateSignature(r.InKey); *sig != *ff.Signature {
			return nil, newError("wrong signature")
		}

		// in UDP, packet order is not guaranteed. Therefore, we accept frames
		// with a timestamp within 10 seconds with respect to the previous
		if r.curReadSignatureTime > 0 &&
			ff.SignatureTimestamp < (r.curReadSignatureTime-(10*100000)) {
			return nil, newError("signature timestamp is too old")
		}

		if ff.SignatureTimestamp > r.curReadSignatureTime {
			r.curReadSignatureTime = ff.SignatureTimestamp
		}
	}

	// decode message if in dialect and validate checksum
	if r.DialectRW != nil {
		if mp := r.DialectRW.GetMessage(f.GetMessage().GetID()); mp != nil {
			if sum := f.GenerateChecksum(mp.CRCExtra()); sum != f.GetChecksum() {
				return nil, newError("wrong checksum, expected %.4x, got %.4x, message id is %d",
					sum, f.GetChecksum(), f.GetMessage().GetID())
			}

			_, isV2 := f.(*V2Frame)
			msg, err := mp.Read(f.GetMessage().(*message.MessageRaw), isV2)
			if err != nil {
				return nil, newError("unable to decode message: %s", err.Error())
			}

			switch ff := f.(type) {
			case *V1Frame:
				ff.Message = msg
			case *V2Frame:
				ff.Message = msg
			}
		}
	}

	return f, nil
}
