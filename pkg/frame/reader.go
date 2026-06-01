package frame

import (
	"bufio"
	"crypto/subtle"
	"fmt"
	"reflect"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
	"github.com/bluenviron/gomavlib/v4/pkg/message"
)

const (
	bufferSize = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
)

func hasStringFields(msg message.Message) bool {
	typ := reflect.TypeOf(msg).Elem()

	for i := range typ.NumField() {
		if typ.Field(i).Type == reflect.TypeFor[string]() {
			return true
		}
	}

	return false
}

func hasEmptyBytes(buf []byte) bool {
	return len(buf) > 1 && buf[len(buf)-1] == 0x00
}

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e ReadError) Error() string {
	return e.str
}

func newError(format string, args ...any) ReadError {
	return ReadError{
		str: fmt.Sprintf(format, args...),
	}
}

// Reader is a Frame reader.
type Reader struct {
	// underlying byte reader.
	BufByteReader *bufio.Reader

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

		if sig := ff.GenerateSignature(r.InKey); subtle.ConstantTimeCompare(sig[:], ff.Signature[:]) != 1 {
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
			rawMessage := f.GetMessage().(*message.MessageRaw)

			var msg message.Message
			msg, err = mp.Read(rawMessage, isV2)
			if err != nil {
				return nil, newError("unable to decode message: %s", err.Error())
			}

			// some libraries generate non-standard messages, in particular:
			// - messages with junk after string termination
			// - v2 messages with trailing empty bytes
			// The specification says that we must support these messages (and we are)
			// but there might be troubles when re-encoding them, since checksum is different.
			// re-compute the checksum.
			if hasStringFields(msg) || (isV2 && hasEmptyBytes(rawMessage.Payload)) {
				raw2 := mp.Write(msg, isV2)
				switch f := f.(type) {
				case *V1Frame:
					f.Message = raw2
					f.Checksum = f.GenerateChecksum(mp.CRCExtra())
				case *V2Frame:
					f.Message = raw2
					f.Checksum = f.GenerateChecksum(mp.CRCExtra())
				}
			}

			switch f := f.(type) {
			case *V1Frame:
				f.Message = msg
			case *V2Frame:
				f.Message = msg
			}
		}
	}

	return f, nil
}
