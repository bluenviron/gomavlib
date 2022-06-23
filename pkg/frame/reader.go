package frame

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/message"
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

func (e *ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) *ReadError {
	return &ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

// ReaderConf is the configuration of a Reader.
type ReaderConf struct {
	// the underlying bytes reader.
	Reader io.Reader

	// (optional) the dialect which contains the messages that will be read.
	// If not provided, messages are decoded into the MessageRaw struct.
	DialectDE *dialect.ReadWriter

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *V2Key
}

// Reader is a Frame reader.
type Reader struct {
	conf                 ReaderConf
	br                   *bufio.Reader
	curReadSignatureTime uint64
}

// NewReader allocates a Reader.
func NewReader(conf ReaderConf) (*Reader, error) {
	if conf.Reader == nil {
		return nil, fmt.Errorf("Reader not provided")
	}

	return &Reader{
		conf: conf,
		br:   bufio.NewReaderSize(conf.Reader, bufferSize),
	}, nil
}

// Read reads a Frame from the reader.
// It must not be called by multiple routines in parallel.
func (r *Reader) Read() (Frame, error) {
	magicByte, err := r.br.ReadByte()
	if err != nil {
		return nil, err
	}

	f, err := func() (Frame, error) {
		switch magicByte {
		case V1MagicByte:
			return &V1Frame{}, nil

		case V2MagicByte:
			return &V2Frame{}, nil
		}

		return nil, newError("invalid magic byte: %x", magicByte)
	}()
	if err != nil {
		return nil, err
	}

	err = f.decode(r.br)
	if err != nil {
		return nil, newError(err.Error())
	}

	if r.conf.InKey != nil {
		ff, ok := f.(*V2Frame)
		if !ok {
			return nil, newError("signature required but packet is not v2")
		}

		if sig := ff.genSignature(r.conf.InKey); *sig != *ff.Signature {
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
	if r.conf.DialectDE != nil {
		if mp, ok := r.conf.DialectDE.MessageDEs[f.GetMessage().GetID()]; ok {
			if sum := f.genChecksum(r.conf.DialectDE.MessageDEs[f.GetMessage().GetID()].CRCExtra()); sum != f.getChecksum() {
				return nil, newError("wrong checksum (expected %.4x, got %.4x, id=%d)",
					sum, f.getChecksum(), f.GetMessage().GetID())
			}

			_, isV2 := f.(*V2Frame)
			msg, err := mp.Read(f.GetMessage().(*message.MessageRaw).Payload, isV2)
			if err != nil {
				return nil, newError(err.Error())
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
