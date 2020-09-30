package gomavlib

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/dialect"
	"github.com/aler9/gomavlib/frame"
	"github.com/aler9/gomavlib/msg"
	"github.com/aler9/gomavlib/x25"
)

func uint24Decode(in []byte) uint32 {
	return uint32(in[2])<<16 | uint32(in[1])<<8 | uint32(in[0])
}

func uint24Encode(buf []byte, in uint32) []byte {
	buf[0] = byte(in)
	buf[1] = byte(in >> 8)
	buf[2] = byte(in >> 16)
	return buf[:3]
}

func uint48Decode(in []byte) uint64 {
	return uint64(in[5])<<40 | uint64(in[4])<<32 | uint64(in[3])<<24 |
		uint64(in[2])<<16 | uint64(in[1])<<8 | uint64(in[0])
}

func uint48Encode(buf []byte, in uint64) []byte {
	buf[0] = byte(in)
	buf[1] = byte(in >> 8)
	buf[2] = byte(in >> 16)
	buf[3] = byte(in >> 24)
	buf[4] = byte(in >> 32)
	buf[5] = byte(in >> 40)
	return buf[:6]
}

// Key is a key able to sign and validate V2 frames.
type Key [32]byte

// NewKey allocates a Key.
func NewKey(in []byte) *Key {
	key := new(Key)
	copy(key[:], in)
	return key
}

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)

// Version is a Mavlink version.
type Version int

const (
	// V1 is Mavlink 1.0
	V1 Version = 1

	// V2 is Mavlink 2.0
	V2 Version = 2
)

// String implements fmt.Stringer and returns the version label.
func (v Version) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}

// ParserError is the error returned in case of non-fatal parsing errors.
type ParserError struct {
	str string
}

func (e *ParserError) Error() string {
	return e.str
}

func newParserError(format string, args ...interface{}) *ParserError {
	return &ParserError{
		str: fmt.Sprintf(format, args...),
	}
}

// ParserConf configures a Parser.
type ParserConf struct {
	// the reader from which frames will be read.
	Reader io.Reader
	// the writer to which frames will be written.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	Dialect *dialect.Dialect

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemId byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentId byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkId byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *Key
}

// Parser is a low-level Mavlink encoder and decoder that works with a Reader and a Writer.
type Parser struct {
	conf                 ParserConf
	dialectDE            *dialect.DecEncoder
	readBuffer           *bufio.Reader
	writeBuffer          []byte
	curWriteSequenceId   byte
	curReadSignatureTime uint64
}

// NewParser allocates a Parser, a low level frame encoder and decoder.
// See ParserConf for the options.
func NewParser(conf ParserConf) (*Parser, error) {
	if conf.Reader == nil {
		return nil, fmt.Errorf("Reader not provided")
	}
	if conf.Writer == nil {
		return nil, fmt.Errorf("Writer not provided")
	}

	if conf.OutVersion == 0 {
		return nil, fmt.Errorf("OutVersion not provided")
	}
	if conf.OutSystemId < 1 {
		return nil, fmt.Errorf("SystemId must be >= 1")
	}
	if conf.OutComponentId < 1 {
		conf.OutComponentId = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	dialectDE, err := func() (*dialect.DecEncoder, error) {
		if conf.Dialect == nil {
			return nil, nil
		}
		return dialect.NewDecEncoder(conf.Dialect)
	}()
	if err != nil {
		return nil, err
	}

	return &Parser{
		conf:        conf,
		dialectDE:   dialectDE,
		readBuffer:  bufio.NewReaderSize(conf.Reader, netBufferSize),
		writeBuffer: make([]byte, 0, netBufferSize),
	}, nil
}

func (p *Parser) checksum(f frame.Frame) uint16 {
	msg := f.GetMessage().(*msg.MessageRaw)
	h := x25.New()

	// the checksum covers the whole message, excluding magic byte, checksum and signature
	switch ff := f.(type) {
	case *frame.V1Frame:
		h.Write([]byte{byte(len(msg.Content))})
		h.Write([]byte{ff.SequenceId})
		h.Write([]byte{ff.SystemId})
		h.Write([]byte{ff.ComponentId})
		h.Write([]byte{byte(msg.Id)})
		h.Write(msg.Content)

	case *frame.V2Frame:
		buf := make([]byte, 3)
		h.Write([]byte{byte(len(msg.Content))})
		h.Write([]byte{ff.IncompatibilityFlag})
		h.Write([]byte{ff.CompatibilityFlag})
		h.Write([]byte{ff.SequenceId})
		h.Write([]byte{ff.SystemId})
		h.Write([]byte{ff.ComponentId})
		h.Write(uint24Encode(buf, msg.Id))
		h.Write(msg.Content)
	}

	// CRC_EXTRA byte is added at the end of the data
	h.Write([]byte{p.dialectDE.MessageDEs[msg.GetId()].CRCExtra()})

	return h.Sum16()
}

func (p *Parser) signature(ff *frame.V2Frame, key *Key) *frame.V2Signature {
	msg := ff.GetMessage().(*msg.MessageRaw)
	h := sha256.New()

	// secret key
	h.Write(key[:])

	// the signature covers the whole message, excluding the signature itself
	buf := make([]byte, 6)
	h.Write([]byte{frame.V2MagicByte})
	h.Write([]byte{byte(len(msg.Content))})
	h.Write([]byte{ff.IncompatibilityFlag})
	h.Write([]byte{ff.CompatibilityFlag})
	h.Write([]byte{ff.SequenceId})
	h.Write([]byte{ff.SystemId})
	h.Write([]byte{ff.ComponentId})
	h.Write(uint24Encode(buf, ff.Message.GetId()))
	h.Write(msg.Content)
	binary.LittleEndian.PutUint16(buf, ff.Checksum)
	h.Write(buf[:2])
	h.Write([]byte{ff.SignatureLinkId})
	h.Write(uint48Encode(buf, ff.SignatureTimestamp))

	sig := new(frame.V2Signature)
	copy(sig[:], h.Sum(nil)[:6])
	return sig
}

// Read reads a Frame from the reader. It must not be called
// by multiple routines in parallel.
func (p *Parser) Read() (frame.Frame, error) {
	magicByte, err := p.readBuffer.ReadByte()
	if err != nil {
		return nil, err
	}

	f, err := func() (frame.Frame, error) {
		switch magicByte {
		case frame.V1MagicByte:
			return &frame.V1Frame{}, nil

		case frame.V2MagicByte:
			return &frame.V2Frame{}, nil
		}

		return nil, newParserError("invalid magic byte: %x", magicByte)
	}()
	if err != nil {
		return nil, err
	}

	err = f.Decode(p.readBuffer)
	if err != nil {
		return nil, newParserError(err.Error())
	}

	if p.conf.InKey != nil {
		ff, ok := f.(*frame.V2Frame)
		if !ok {
			return nil, newParserError("signature required but packet is not v2")
		}

		if sig := p.signature(ff, p.conf.InKey); *sig != *ff.Signature {
			return nil, newParserError("wrong signature")
		}

		// in UDP, packet order is not guaranteed. Therefore, we accept frames
		// with a timestamp within 10 seconds with respect to the previous frame.
		if p.curReadSignatureTime > 0 &&
			ff.SignatureTimestamp < (p.curReadSignatureTime-(10*100000)) {
			return nil, newParserError("signature timestamp is too old")
		}

		if ff.SignatureTimestamp > p.curReadSignatureTime {
			p.curReadSignatureTime = ff.SignatureTimestamp
		}
	}

	// decode message if in dialect and validate checksum
	if p.conf.Dialect != nil {
		if mp, ok := p.dialectDE.MessageDEs[f.GetMessage().GetId()]; ok {
			if sum := p.checksum(f); sum != f.GetChecksum() {
				return nil, newParserError("wrong checksum (expected %.4x, got %.4x, id=%d)",
					sum, f.GetChecksum(), f.GetMessage().GetId())
			}

			_, isV2 := f.(*frame.V2Frame)
			msg, err := mp.Decode(f.GetMessage().(*msg.MessageRaw).Content, isV2)
			if err != nil {
				return nil, newParserError(err.Error())
			}

			switch ff := f.(type) {
			case *frame.V1Frame:
				ff.Message = msg
			case *frame.V2Frame:
				ff.Message = msg
			}
		}
	}

	return f, nil
}

// WriteMessage writes a Message into the writer.
// It must not be called by multiple routines in parallel.
func (p *Parser) WriteMessage(message msg.Message) error {
	var f frame.Frame
	if p.conf.OutVersion == V1 {
		f = &frame.V1Frame{Message: message}
	} else {
		f = &frame.V2Frame{Message: message}
	}
	return p.writeFrameAndFill(f)
}

func (p *Parser) writeFrameAndFill(f frame.Frame) error {
	if f.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// do not touch the original frame, but work with a separate object
	// in such way that the frame can be encoded by other parsers in parallel
	safeFrame := f.Clone()

	// fill SequenceId, SystemId, ComponentId
	switch ff := safeFrame.(type) {
	case *frame.V1Frame:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	case *frame.V2Frame:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	}
	p.curWriteSequenceId++

	// fill CompatibilityFlag, IncompatibilityFlag if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok {
		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0

		if p.conf.OutKey != nil {
			ff.IncompatibilityFlag |= frame.V2FlagSigned
		}
	}

	// encode message if it is not already encoded
	if _, ok := safeFrame.GetMessage().(*msg.MessageRaw); !ok {
		if p.conf.Dialect == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.dialectDE.MessageDEs[safeFrame.GetMessage().GetId()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := safeFrame.(*frame.V2Frame)
		byt, err := mp.Encode(safeFrame.GetMessage(), isV2)
		if err != nil {
			return err
		}

		msgRaw := &msg.MessageRaw{safeFrame.GetMessage().GetId(), byt}
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Message = msgRaw
		case *frame.V2Frame:
			ff.Message = msgRaw
		}

		// fill checksum
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Checksum = p.checksum(ff)
		case *frame.V2Frame:
			ff.Checksum = p.checksum(ff)
		}
	}

	// fill SignatureLinkId, SignatureTimestamp, Signature if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok && p.conf.OutKey != nil {
		ff.SignatureLinkId = p.conf.OutSignatureLinkId
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = p.signature(ff, p.conf.OutKey)
	}

	return p.WriteFrame(safeFrame)
}

// WriteFrame writes a Frame into the writer.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (p *Parser) WriteFrame(f frame.Frame) error {
	m := f.GetMessage()
	if m == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := m.(*msg.MessageRaw); !ok {
		if p.conf.Dialect == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.dialectDE.MessageDEs[m.GetId()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := f.(*frame.V2Frame)
		byt, err := mp.Encode(m, isV2)
		if err != nil {
			return err
		}

		// do not touch frame.Message
		// in such way that the frame can be encoded by other parsers in parallel
		m = &msg.MessageRaw{m.GetId(), byt}
	}

	buf, err := f.Encode(p.writeBuffer, m.(*msg.MessageRaw).Content)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = p.conf.Writer.Write(buf)
	return err
}
