package gomavlib

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
)

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

// ParserReader is the interface that must be implemented by readers passed to Read()
type ParserReader interface {
	io.ByteReader
	io.Reader
}

// ParserConf configures a Parser.
type ParserConf struct {
	// the reader from which frames will be read.
	//Reader io.Reader
	// the writer to which frames will be written.
	//Writer io.Writer

	// contains the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect []Message

	// (optional) the secret key used to verify incoming frames.
	// Non signed frames are discarded. This feature requires Mavlink v2.
	SignatureInKey *FrameSignatureKey

	// (optional) disables checksum validation of incoming frames.
	// Not recommended, useful only for debugging purposes.
	ChecksumDisable bool
}

// Parser is a low-level Mavlink encoder and decoder that works with a Reader and a Writer.
type Parser struct {
	conf           ParserConf
	parserMessages map[uint32]*parserMessage
}

// NewParser allocates a Parser, a low level frame encoder and decoder.
// See Parser for the options.
func NewParser(conf ParserConf) (*Parser, error) {
	p := &Parser{
		conf:           conf,
		parserMessages: make(map[uint32]*parserMessage),
	}

	// generate message parsers
	for _, msg := range conf.Dialect {
		mp, err := newParserMessage(msg)
		if err != nil {
			return nil, fmt.Errorf("message %T: %s", msg, err)
		}
		p.parserMessages[msg.GetId()] = mp
	}

	return p, nil
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
	h.Write([]byte{p.parserMessages[msg.GetId()].crcExtra})

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

// Read returns the first Frame obtained from reader.
func (p *Parser) Read(reader ParserReader) (Frame, error) {
	magicByte, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	var f Frame
	switch magicByte {
	case v1MagicByte:
		ff := &FrameV1{}
		f = ff

		// header
		msgLen, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SequenceId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SystemId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.ComponentId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		msgId, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}

		// message
		msgContent := make([]byte, msgLen)
		_, err = io.ReadFull(reader, msgContent)
		if err != nil {
			return nil, err
		}
		ff.Message = &MessageRaw{
			Id:      uint32(msgId),
			Content: msgContent,
		}

		// checksum
		err = binary.Read(reader, binary.LittleEndian, &ff.Checksum)
		if err != nil {
			return nil, err
		}

	case v2MagicByte:
		ff := &FrameV2{}
		f = ff

		// header
		msgLen, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.IncompatibilityFlag, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.CompatibilityFlag, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SequenceId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.SystemId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		ff.ComponentId, err = reader.ReadByte()
		if err != nil {
			return nil, err
		}
		var msgId uint32
		err = readLittleEndianUint24(reader, &msgId)
		if err != nil {
			return nil, err
		}

		// discard frame if incompatibility flag is not understood, as in recommendations
		if ff.IncompatibilityFlag != 0 && ff.IncompatibilityFlag != flagSigned {
			return nil, newParserError("unknown incompatibility flag (%d)", ff.IncompatibilityFlag)
		}

		// message
		msgContent := make([]byte, msgLen)
		_, err = io.ReadFull(reader, msgContent)
		if err != nil {
			return nil, err
		}
		ff.Message = &MessageRaw{
			Id:      msgId,
			Content: msgContent,
		}

		// checksum
		err = binary.Read(reader, binary.LittleEndian, &ff.Checksum)
		if err != nil {
			return nil, err
		}

		// signature
		if ff.IsSigned() {
			ff.SignatureLinkId, err = reader.ReadByte()
			if err != nil {
				return nil, err
			}
			err = readLittleEndianUint48(reader, &ff.SignatureTimestamp)
			if err != nil {
				return nil, err
			}
			ff.Signature = new(FrameSignature)
			_, err = io.ReadFull(reader, ff.Signature[:])
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
	if mp, ok := p.parserMessages[f.GetMessage().GetId()]; ok {
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

	return f, nil
}

// Write writes a Frame into writer.
func (p *Parser) Write(writer io.Writer, f Frame, fillChecksum bool, fillSignatureKey *FrameSignatureKey) error {
	// encode message if not already encoded and in dialect
	if _, ok := f.GetMessage().(*MessageRaw); ok == false {
		if mp, ok := p.parserMessages[f.GetMessage().GetId()]; ok {
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
			if ff, ok := f.(*FrameV2); ok && fillSignatureKey != nil {
				ff.IncompatibilityFlag |= flagSigned
			}

			if fillChecksum == true {
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

	msgContent := f.GetMessage().(*MessageRaw).Content

	var buf []byte
	switch ff := f.(type) {
	case *FrameV1:
		if ff.Message.GetId() > 0xFF {
			return fmt.Errorf("cannot send a message with an id > 0xFF and a V1 frame")
		}

		bufferLen := 6 + len(msgContent) + 2
		buf = make([]byte, bufferLen)

		offset := 0
		buf[offset] = v1MagicByte
		offset += 1
		buf[offset] = byte(len(msgContent))
		offset += 1
		buf[offset] = ff.SequenceId
		offset += 1
		buf[offset] = ff.SystemId
		offset += 1
		buf[offset] = ff.ComponentId
		offset += 1

		buf[offset] = byte(ff.Message.GetId())
		offset += 1

		copy(buf[offset:], msgContent)
		offset += len(msgContent)

		binary.LittleEndian.PutUint16(buf[offset:offset+2], ff.Checksum)

	case *FrameV2:
		if fillSignatureKey != nil {
			ff.Signature = p.Signature(ff, fillSignatureKey)
		}

		bufferLen := 10 + len(msgContent) + 2
		if ff.IsSigned() {
			bufferLen += 13
		}
		buf = make([]byte, bufferLen)

		offset := 0
		buf[offset] = v2MagicByte
		offset += 1
		buf[offset] = byte(len(msgContent))
		offset += 1
		buf[offset] = ff.IncompatibilityFlag
		offset += 1
		buf[offset] = ff.CompatibilityFlag
		offset += 1
		buf[offset] = ff.SequenceId
		offset += 1
		buf[offset] = ff.SystemId
		offset += 1
		buf[offset] = ff.ComponentId
		offset += 1

		copy(buf[7:10], uint24Encode(ff.Message.GetId()))
		offset += 3

		copy(buf[offset:], msgContent)
		offset += len(msgContent)

		binary.LittleEndian.PutUint16(buf[offset:offset+2], ff.Checksum)
		offset += 2

		if ff.IsSigned() {
			buf[offset] = ff.SignatureLinkId
			offset += 1
			copy(buf[offset:offset+6], uint48Encode(ff.SignatureTimestamp))
			offset += 6
			copy(buf[offset:offset+6], ff.Signature[:])
		}
	}

	// we do not check n
	// since io.Writer is not allowed to return n < len(buf) without throwing an error
	_, err := writer.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
