package gomavlib

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)

// ParserError is the error returned in case of non-fatal parsing errors
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

	// contains the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect *Dialect

	// these are used to identify this node in the network.
	// They are added to every outgoing frame.
	SystemId    byte
	ComponentId byte

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires Mavlink v2.
	SignatureInKey *FrameSignatureKey
	// (optional) the value to insert into the signature link id
	SignatureLinkId byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires Mavlink v2.
	SignatureOutKey *FrameSignatureKey

	// (optional) disables checksum validation of incoming frames.
	// Not recommended, useful only for debugging purposes.
	ChecksumDisable bool
}

// Parser is a low-level Mavlink encoder and decoder that works with a Reader and a Writer.
type Parser struct {
	conf           ParserConf
	readBuffer     *bufio.Reader
	writeBuffer    []byte
	nextSequenceId byte
}

// NewParser allocates a Parser, a low level frame encoder and decoder.
// See Parser for the options.
func NewParser(conf ParserConf) *Parser {
	p := &Parser{
		conf:        conf,
		readBuffer:  bufio.NewReaderSize(conf.Reader, netBufferSize),
		writeBuffer: make([]byte, 0, netBufferSize),
	}
	return p
}

// Checksum computes the checksum of a given frame.
func (p *Parser) Checksum(f Frame) uint16 {
	msg := f.GetMessage().(*MessageRaw)
	h := NewX25()

	// the checksum covers the whole message, excluding magic byte, checksum and signature
	switch ff := f.(type) {
	case *FrameV1:
		h.Write([]byte{byte(len(msg.Content))})
		h.Write([]byte{ff.SequenceId})
		h.Write([]byte{ff.SystemId})
		h.Write([]byte{ff.ComponentId})
		h.Write([]byte{byte(msg.Id)})
		h.Write(msg.Content)

	case *FrameV2:
		h.Write([]byte{byte(len(msg.Content))})
		h.Write([]byte{ff.IncompatibilityFlag})
		h.Write([]byte{ff.CompatibilityFlag})
		h.Write([]byte{ff.SequenceId})
		h.Write([]byte{ff.SystemId})
		h.Write([]byte{ff.ComponentId})
		h.Write(uint24Encode(msg.Id))
		h.Write(msg.Content)
	}

	// CRC_EXTRA byte is added at the end of the data
	h.Write([]byte{p.conf.Dialect.messages[msg.GetId()].crcExtra})

	return h.Sum16()
}

// Signature computes the signature of a given frame with the given key.
func (p *Parser) Signature(ff *FrameV2, key *FrameSignatureKey) *FrameSignature {
	msg := ff.GetMessage().(*MessageRaw)
	h := sha256.New()

	// secret key
	h.Write(key[:])

	// the signature covers the whole message, excluding the signature itself
	h.Write([]byte{v2MagicByte})
	h.Write([]byte{byte(len(msg.Content))})
	h.Write([]byte{ff.IncompatibilityFlag})
	h.Write([]byte{ff.CompatibilityFlag})
	h.Write([]byte{ff.SequenceId})
	h.Write([]byte{ff.SystemId})
	h.Write([]byte{ff.ComponentId})
	h.Write(uint24Encode(ff.Message.GetId()))
	h.Write(msg.Content)
	binary.Write(h, binary.LittleEndian, ff.Checksum)
	h.Write([]byte{ff.SignatureLinkId})
	h.Write(uint48Encode(ff.SignatureTimestamp))

	sig := new(FrameSignature)
	copy(sig[:], h.Sum(nil)[:6])
	return sig
}

// Read returns the first Frame parsed from the reader. It must not be called
// by multiple routines in parallel.
func (p *Parser) Read() (Frame, error) {
	magicByte, err := p.readBuffer.ReadByte()
	if err != nil {
		return nil, err
	}

	var f Frame
	switch magicByte {
	case v1MagicByte:
		ff := &FrameV1{}
		f = ff

		// header
		msgLen, err := p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SequenceId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SystemId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.ComponentId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		msgId, err := p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}

		// message
		msgContent := make([]byte, msgLen)
		_, err = io.ReadFull(p.readBuffer, msgContent)
		if err != nil {
			return nil, err
		}
		ff.Message = &MessageRaw{
			Id:      uint32(msgId),
			Content: msgContent,
		}

		// checksum
		err = binary.Read(p.readBuffer, binary.LittleEndian, &ff.Checksum)
		if err != nil {
			return nil, err
		}

	case v2MagicByte:
		ff := &FrameV2{}
		f = ff

		// header
		msgLen, err := p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.IncompatibilityFlag, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.CompatibilityFlag, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SequenceId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SystemId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.ComponentId, err = p.readBuffer.ReadByte()
		if err != nil {
			return nil, err
		}
		var msgId uint32
		err = readLittleEndianUint24(p.readBuffer, &msgId)
		if err != nil {
			return nil, err
		}

		// discard frame if incompatibility flag is not understood, as in recommendations
		if ff.IncompatibilityFlag != 0 && ff.IncompatibilityFlag != flagSigned {
			return nil, newParserError("unknown incompatibility flag (%d)", ff.IncompatibilityFlag)
		}

		// message
		msgContent := make([]byte, msgLen)
		_, err = io.ReadFull(p.readBuffer, msgContent)
		if err != nil {
			return nil, err
		}
		ff.Message = &MessageRaw{
			Id:      msgId,
			Content: msgContent,
		}

		// checksum
		err = binary.Read(p.readBuffer, binary.LittleEndian, &ff.Checksum)
		if err != nil {
			return nil, err
		}

		// signature
		if ff.IsSigned() {
			ff.SignatureLinkId, err = p.readBuffer.ReadByte()
			if err != nil {
				return nil, err
			}
			err = readLittleEndianUint48(p.readBuffer, &ff.SignatureTimestamp)
			if err != nil {
				return nil, err
			}
			ff.Signature = new(FrameSignature)
			_, err = io.ReadFull(p.readBuffer, ff.Signature[:])
			if err != nil {
				return nil, err
			}
		}

	default:
		return nil, newParserError("unrecognized magic byte: %x", magicByte)
	}

	if p.conf.SignatureInKey != nil {
		ff, ok := f.(*FrameV2)
		if ok == false {
			return nil, newParserError("signature required but packet is not v2")
		}

		if sig := p.Signature(ff, p.conf.SignatureInKey); *sig != *ff.Signature {
			return nil, newParserError("wrong signature")
		}
	}

	// decode message if in dialect and validate checksum
	if p.conf.Dialect != nil {
		if mp, ok := p.conf.Dialect.messages[f.GetMessage().GetId()]; ok {
			if p.conf.ChecksumDisable == false {
				if sum := p.Checksum(f); sum != f.GetChecksum() {
					return nil, newParserError("wrong checksum (expected %.4x, got %.4x, id=%d)",
						sum, f.GetChecksum(), f.GetMessage().GetId())
				}
			}

			_, isFrameV2 := f.(*FrameV2)
			msg, err := mp.decode(f.GetMessage().(*MessageRaw).Content, isFrameV2)
			if err != nil {
				return nil, newParserError(err.Error())
			}

			switch ff := f.(type) {
			case *FrameV1:
				ff.Message = msg
			case *FrameV2:
				ff.Message = msg
			}
		}
	}

	return f, nil
}

// Write writes a Frame into the writer. It must not be called by multiple
// routines in parallel.
// if route is false, the following fields will be filled:
// - sequence id
// - system id
// - component id
// - checksum
// - signature link id
// - signature timestamp
// - signature
// if route is true, the frame will be left untouched.
func (p *Parser) Write(f Frame, route bool) error {
	if route == false {
		switch ff := f.(type) {
		case *FrameV1:
			ff.SequenceId = p.nextSequenceId
			ff.SystemId = p.conf.SystemId
			ff.ComponentId = p.conf.ComponentId
		case *FrameV2:
			ff.SequenceId = p.nextSequenceId
			ff.SystemId = p.conf.SystemId
			ff.ComponentId = p.conf.ComponentId
		}
		p.nextSequenceId++
	}

	// encode message if not already encoded and in dialect
	if p.conf.Dialect != nil {
		if _, ok := f.GetMessage().(*MessageRaw); ok == false {
			if mp, ok := p.conf.Dialect.messages[f.GetMessage().GetId()]; ok {
				_, isFrameV2 := f.(*FrameV2)
				byt, err := mp.encode(f.GetMessage(), isFrameV2)
				if err != nil {
					return err
				}

				switch ff := f.(type) {
				case *FrameV1:
					ff.Message = &MessageRaw{f.GetMessage().GetId(), byt}
				case *FrameV2:
					ff.Message = &MessageRaw{f.GetMessage().GetId(), byt}
				}

				// if frame is going to be signed, set incompatibility flag
				// before computing checksum
				if ff, ok := f.(*FrameV2); ok && route == false && p.conf.SignatureOutKey != nil {
					ff.IncompatibilityFlag |= flagSigned
				}

				if route == false {
					check := p.Checksum(f)
					switch ff := f.(type) {
					case *FrameV1:
						ff.Checksum = check
					case *FrameV2:
						ff.Checksum = check
					}
				}
			}
		}
	}

	msgContent := f.GetMessage().(*MessageRaw).Content

	offset := 0
	switch ff := f.(type) {
	case *FrameV1:
		if ff.Message.GetId() > 0xFF {
			return fmt.Errorf("cannot send a message with an id > 0xFF and a V1 frame")
		}

		bufferLen := 6 + len(msgContent) + 2
		p.writeBuffer = p.writeBuffer[:bufferLen]

		// header
		p.writeBuffer[offset] = v1MagicByte
		offset += 1
		p.writeBuffer[offset] = byte(len(msgContent))
		offset += 1
		p.writeBuffer[offset] = ff.SequenceId
		offset += 1
		p.writeBuffer[offset] = ff.SystemId
		offset += 1
		p.writeBuffer[offset] = ff.ComponentId
		offset += 1
		p.writeBuffer[offset] = byte(ff.Message.GetId())
		offset += 1

		// message
		copy(p.writeBuffer[offset:], msgContent)
		offset += len(msgContent)

		binary.LittleEndian.PutUint16(p.writeBuffer[offset:offset+2], ff.Checksum)

	case *FrameV2:
		if route == false && p.conf.SignatureOutKey != nil {
			ff.SignatureLinkId = p.conf.SignatureLinkId
			// Timestamp in 10 microsecond units since 1st January 2015 GMT time
			ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
			ff.Signature = p.Signature(ff, p.conf.SignatureOutKey)
		}

		bufferLen := 10 + len(msgContent) + 2
		if ff.IsSigned() {
			bufferLen += 13
		}
		p.writeBuffer = p.writeBuffer[:bufferLen]

		// header
		p.writeBuffer[offset] = v2MagicByte
		offset += 1
		p.writeBuffer[offset] = byte(len(msgContent))
		offset += 1
		p.writeBuffer[offset] = ff.IncompatibilityFlag
		offset += 1
		p.writeBuffer[offset] = ff.CompatibilityFlag
		offset += 1
		p.writeBuffer[offset] = ff.SequenceId
		offset += 1
		p.writeBuffer[offset] = ff.SystemId
		offset += 1
		p.writeBuffer[offset] = ff.ComponentId
		offset += 1
		copy(p.writeBuffer[offset:offset+3], uint24Encode(ff.Message.GetId()))
		offset += 3

		// message
		copy(p.writeBuffer[offset:], msgContent)
		offset += len(msgContent)

		binary.LittleEndian.PutUint16(p.writeBuffer[offset:offset+2], ff.Checksum)
		offset += 2

		if ff.IsSigned() {
			p.writeBuffer[offset] = ff.SignatureLinkId
			offset += 1
			copy(p.writeBuffer[offset:offset+6], uint48Encode(ff.SignatureTimestamp))
			offset += 6
			copy(p.writeBuffer[offset:offset+6], ff.Signature[:])
		}
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err := p.conf.Writer.Write(p.writeBuffer)
	return err
}
