package gomavlib

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// ParserConf configures a Parser.
type ParserConf struct {
	// Dialect contains the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect []Message
}

// Parser is a low-level Mavlink encoder and decoder that works with byte
// slices.
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

// Decode converts a byte buffer to a Frame.
func (p *Parser) Decode(buf []byte, validateChecksum bool, validateSignatureKey *FrameSignatureKey) (Frame, error) {
	// require at least magic byte, message length and incompatibility flag (if v2)
	if len(buf) < 3 {
		return nil, fmt.Errorf("insufficient packet length")
	}

	var f Frame
	switch buf[0] {
	case v1MagicByte:
		ff := &FrameV1{}
		f = ff

		bufferLen := 6 + int(buf[1]) + 2
		if len(buf) != bufferLen {
			return nil, fmt.Errorf("wrong packet length (got %d, expected %d)", len(buf), bufferLen)
		}

		offset := 2
		ff.SequenceId = buf[offset]
		offset += 1
		ff.SystemId = buf[offset]
		offset += 1
		ff.ComponentId = buf[offset]
		offset += 1

		id := uint32(buf[offset])
		offset += 1
		content := buf[offset : offset+int(buf[1])]
		offset += int(buf[1])

		ff.Message = &MessageRaw{
			Id:      id,
			Content: content,
		}

		ff.Checksum = binary.LittleEndian.Uint16(buf[offset : offset+2])

	case v2MagicByte:
		ff := &FrameV2{}
		f = ff

		offset := 2
		ff.IncompatibilityFlag = buf[offset]
		offset += 1

		// discard frame if incompatibility flag is not understood, as in recommendations
		if ff.IncompatibilityFlag != 0 && ff.IncompatibilityFlag != 0x01 {
			return nil, fmt.Errorf("unknown incompatibility flag (%d)", ff.IncompatibilityFlag)
		}

		bufferLen := 10 + int(buf[1]) + 2
		if ff.IsSigned() {
			bufferLen += 13
		}
		if bufferLen != len(buf) {
			return nil, fmt.Errorf("wrong packet length (got %d, expected %d)", len(buf), bufferLen)
		}

		ff.CompatibilityFlag = buf[offset]
		offset += 1
		ff.SequenceId = buf[offset]
		offset += 1
		ff.SystemId = buf[offset]
		offset += 1
		ff.ComponentId = buf[offset]
		offset += 1

		id := uint24Decode(buf[offset : offset+3])
		offset += 3
		content := buf[offset : offset+int(buf[1])]
		offset += int(buf[1])

		ff.Message = &MessageRaw{
			Id:      id,
			Content: content,
		}

		ff.Checksum = binary.LittleEndian.Uint16(buf[offset : offset+2])
		offset += 2

		if ff.IsSigned() {
			ff.SignatureLinkId = buf[offset]
			offset += 1
			ff.SignatureTimestamp = uint48Decode(buf[offset : offset+6])
			offset += 6
			ff.Signature = new(FrameSignature)
			copy(ff.Signature[:], buf[offset:offset+6])
		}

	default:
		return nil, fmt.Errorf("unrecognized magic byte: %x", buf[0])
	}

	if validateSignatureKey != nil {
		ff, ok := f.(*FrameV2)
		if ok == false {
			return nil, fmt.Errorf("signature required but packet is not v2")
		}

		if sig := p.Signature(ff, validateSignatureKey); *sig != *ff.Signature {
			return nil, fmt.Errorf("wrong signature")
		}
	}

	// decode message if in dialect
	if mp, ok := p.parserMessages[f.GetMessage().GetId()]; ok {

		if validateChecksum == true {
			if sum := p.Checksum(f); sum != f.GetChecksum() {
				return nil, fmt.Errorf("wrong checksum (expected %.4x, got %.4x)", sum, f.GetChecksum())
			}
		}

		_, isFrameV2 := f.(*FrameV2)
		msg, err := mp.decode(f.GetMessage().(*MessageRaw).Content, isFrameV2)
		if err != nil {
			return nil, err
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

// Encode converts a Frame into a bytes buffer.
func (p *Parser) Encode(f Frame, fillChecksum bool, fillSignatureKey *FrameSignatureKey) ([]byte, error) {
	// encode message if not already encoded and in dialect
	if _, ok := f.GetMessage().(*MessageRaw); ok == false {
		if mp, ok := p.parserMessages[f.GetMessage().GetId()]; ok {
			_, isFrameV2 := f.(*FrameV2)
			byt, err := mp.encode(f.GetMessage(), isFrameV2)
			if err != nil {
				return nil, err
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
			return nil, fmt.Errorf("cannot send a message with an id > 0xFF and a V1 frame")
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
	return buf, nil
}