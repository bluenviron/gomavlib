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

// Version allows to set the frame version used to wrap outgoing messages.
type Version int

const (
	// V2 (default) wraps outgoing messages in v2 frames.
	V2 Version = iota + 1
	// V1 wraps outgoing messages in v1 frames.
	V1
)

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

	// (optional) the messages which will be automatically decoded and
	// encoded. If not provided, messages are decoded in the MessageRaw struct.
	Dialect *Dialect

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

	p := &Parser{
		conf:        conf,
		readBuffer:  bufio.NewReaderSize(conf.Reader, _NET_BUFFER_SIZE),
		writeBuffer: make([]byte, 0, _NET_BUFFER_SIZE),
	}
	return p, nil
}

func (p *Parser) checksum(f Frame) uint16 {
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
	h.Write([]byte{p.conf.Dialect.messages[msg.GetId()].crcExtra})

	return h.Sum16()
}

func (p *Parser) signature(ff *FrameV2, key *Key) *Signature {
	msg := ff.GetMessage().(*MessageRaw)
	h := sha256.New()

	// secret key
	h.Write(key[:])

	// the signature covers the whole message, excluding the signature itself
	buf := make([]byte, 6)
	h.Write([]byte{v2MagicByte})
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

	sig := new(Signature)
	copy(sig[:], h.Sum(nil)[:6])
	return sig
}

// Read reads a Frame from the reader. It must not be called
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
		buf, err := p.readBuffer.Peek(5)
		if err != nil {
			return nil, err
		}
		p.readBuffer.Discard(5)
		msgLen := buf[0]
		ff.SequenceId = buf[1]
		ff.SystemId = buf[2]
		ff.ComponentId = buf[3]
		msgId := buf[4]

		// message
		var msgContent []byte
		if msgLen > 0 {
			msgContent = make([]byte, msgLen)
			_, err = io.ReadFull(p.readBuffer, msgContent)
			if err != nil {
				return nil, err
			}
		}
		ff.Message = &MessageRaw{
			Id:      uint32(msgId),
			Content: msgContent,
		}

		// checksum
		buf, err = p.readBuffer.Peek(2)
		if err != nil {
			return nil, err
		}
		p.readBuffer.Discard(2)
		ff.Checksum = binary.LittleEndian.Uint16(buf)

	case v2MagicByte:
		ff := &FrameV2{}
		f = ff

		// header
		buf, err := p.readBuffer.Peek(9)
		if err != nil {
			return nil, err
		}
		p.readBuffer.Discard(9)
		msgLen := buf[0]
		ff.IncompatibilityFlag = buf[1]
		ff.CompatibilityFlag = buf[2]
		ff.SequenceId = buf[3]
		ff.SystemId = buf[4]
		ff.ComponentId = buf[5]
		msgId := uint24Decode(buf[6:])

		// discard frame if incompatibility flag is not understood, as in recommendations
		if ff.IncompatibilityFlag != 0 && ff.IncompatibilityFlag != flagSigned {
			return nil, newParserError("unknown incompatibility flag (%d)", ff.IncompatibilityFlag)
		}

		// message
		var msgContent []byte
		if msgLen > 0 {
			msgContent = make([]byte, msgLen)
			_, err = io.ReadFull(p.readBuffer, msgContent)
			if err != nil {
				return nil, err
			}
		}
		ff.Message = &MessageRaw{
			Id:      msgId,
			Content: msgContent,
		}

		// checksum
		buf, err = p.readBuffer.Peek(2)
		if err != nil {
			return nil, err
		}
		p.readBuffer.Discard(2)
		ff.Checksum = binary.LittleEndian.Uint16(buf)

		// signature
		if ff.IsSigned() {
			buf, err := p.readBuffer.Peek(13)
			if err != nil {
				return nil, err
			}
			p.readBuffer.Discard(13)
			ff.SignatureLinkId = buf[0]
			ff.SignatureTimestamp = uint48Decode(buf[1:])
			ff.Signature = new(Signature)
			copy(ff.Signature[:], buf[7:])
		}

	default:
		return nil, newParserError("unrecognized magic byte: %x", magicByte)
	}

	if p.conf.InKey != nil {
		ff, ok := f.(*FrameV2)
		if ok == false {
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
		if mp, ok := p.conf.Dialect.messages[f.GetMessage().GetId()]; ok {
			if sum := p.checksum(f); sum != f.GetChecksum() {
				return nil, newParserError("wrong checksum (expected %.4x, got %.4x, id=%d)",
					sum, f.GetChecksum(), f.GetMessage().GetId())
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

// WriteMessage writes a Message into the writer.
// It must not be called by multiple routines in parallel.
func (p *Parser) WriteMessage(message Message) error {
	var f Frame
	if p.conf.OutVersion == V1 {
		f = &FrameV1{Message: message}
	} else {
		f = &FrameV2{Message: message}
	}
	return p.writeFrameAndFill(f)
}

func (p *Parser) writeFrameAndFill(frame Frame) error {
	if frame.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// do not touch the original frame, but work with a separate object
	// in such way that the frame can be encoded by other parsers in parallel
	safeFrame := frame.Clone()

	// fill SequenceId, SystemId, ComponentId
	switch ff := safeFrame.(type) {
	case *FrameV1:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	case *FrameV2:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	}
	p.curWriteSequenceId++

	// fill CompatibilityFlag, IncompatibilityFlag if v2
	if ff, ok := safeFrame.(*FrameV2); ok {
		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0

		if p.conf.OutKey != nil {
			ff.IncompatibilityFlag |= flagSigned
		}
	}

	// encode message if it is not already encoded
	if _, ok := safeFrame.GetMessage().(*MessageRaw); !ok {
		if p.conf.Dialect == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.conf.Dialect.messages[safeFrame.GetMessage().GetId()]
		if ok == false {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isFrameV2 := safeFrame.(*FrameV2)
		byt, err := mp.encode(safeFrame.GetMessage(), isFrameV2)
		if err != nil {
			return err
		}

		msgRaw := &MessageRaw{safeFrame.GetMessage().GetId(), byt}
		switch ff := safeFrame.(type) {
		case *FrameV1:
			ff.Message = msgRaw
		case *FrameV2:
			ff.Message = msgRaw
		}

		// fill checksum
		switch ff := safeFrame.(type) {
		case *FrameV1:
			ff.Checksum = p.checksum(ff)
		case *FrameV2:
			ff.Checksum = p.checksum(ff)
		}
	}

	// fill SignatureLinkId, SignatureTimestamp, Signature if v2
	if ff, ok := safeFrame.(*FrameV2); ok && p.conf.OutKey != nil {
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
func (p *Parser) WriteFrame(frame Frame) error {
	if frame.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	msg := frame.GetMessage()
	if _, ok := msg.(*MessageRaw); !ok {
		if p.conf.Dialect == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.conf.Dialect.messages[msg.GetId()]
		if ok == false {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isFrameV2 := frame.(*FrameV2)
		byt, err := mp.encode(msg, isFrameV2)
		if err != nil {
			return err
		}

		// do not touch ff.Message
		// in such way that the frame can be encoded by other parsers in parallel
		msg = &MessageRaw{msg.GetId(), byt}
	}

	msgContent := msg.(*MessageRaw).Content
	msgLen := len(msgContent)

	switch ff := frame.(type) {
	case *FrameV1:
		if ff.Message.GetId() > 0xFF {
			return fmt.Errorf("cannot send a message with an id > 0xFF and a V1 frame")
		}

		bufferLen := 6 + msgLen + 2
		p.writeBuffer = p.writeBuffer[:bufferLen]

		// header
		p.writeBuffer[0] = v1MagicByte
		p.writeBuffer[1] = byte(msgLen)
		p.writeBuffer[2] = ff.SequenceId
		p.writeBuffer[3] = ff.SystemId
		p.writeBuffer[4] = ff.ComponentId
		p.writeBuffer[5] = byte(ff.Message.GetId())

		// message
		if msgLen > 0 {
			copy(p.writeBuffer[6:], msgContent)
		}

		// checksum
		binary.LittleEndian.PutUint16(p.writeBuffer[6+msgLen:], ff.Checksum)

	case *FrameV2:
		bufferLen := 10 + msgLen + 2
		if ff.IsSigned() {
			bufferLen += 13
		}
		p.writeBuffer = p.writeBuffer[:bufferLen]

		// header
		p.writeBuffer[0] = v2MagicByte
		p.writeBuffer[1] = byte(msgLen)
		p.writeBuffer[2] = ff.IncompatibilityFlag
		p.writeBuffer[3] = ff.CompatibilityFlag
		p.writeBuffer[4] = ff.SequenceId
		p.writeBuffer[5] = ff.SystemId
		p.writeBuffer[6] = ff.ComponentId
		uint24Encode(p.writeBuffer[7:], ff.Message.GetId())

		// message
		if msgLen > 0 {
			copy(p.writeBuffer[10:], msgContent)
		}

		// checksum
		binary.LittleEndian.PutUint16(p.writeBuffer[10+msgLen:], ff.Checksum)

		// signature
		if ff.IsSigned() {
			p.writeBuffer[12+msgLen] = ff.SignatureLinkId
			uint48Encode(p.writeBuffer[13+msgLen:], ff.SignatureTimestamp)
			copy(p.writeBuffer[19+msgLen:], ff.Signature[:])
		}
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err := p.conf.Writer.Write(p.writeBuffer)
	return err
}
