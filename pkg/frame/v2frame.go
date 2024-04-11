package frame

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/bluenviron/gomavlib/v2/pkg/message"
	"github.com/bluenviron/gomavlib/v2/pkg/x25"
)

const (
	// V2MagicByte is the magic byte of a V2 frame.
	V2MagicByte = 0xFD

	// V2FlagSigned is the flag of a V2 frame that indicates that the frame is signed.
	V2FlagSigned = 0x01
)

func uint24Decode(in []byte) uint32 {
	return uint32(in[2])<<16 | uint32(in[1])<<8 | uint32(in[0])
}

func uint24Encode(buf []byte, in uint32) {
	buf[0] = byte(in)
	buf[1] = byte(in >> 8)
	buf[2] = byte(in >> 16)
}

func uint48Decode(in []byte) uint64 {
	return uint64(in[5])<<40 | uint64(in[4])<<32 | uint64(in[3])<<24 |
		uint64(in[2])<<16 | uint64(in[1])<<8 | uint64(in[0])
}

func uint48Encode(buf []byte, in uint64) {
	buf[0] = byte(in)
	buf[1] = byte(in >> 8)
	buf[2] = byte(in >> 16)
	buf[3] = byte(in >> 24)
	buf[4] = byte(in >> 32)
	buf[5] = byte(in >> 40)
}

// V2Key is a key able to sign and validate V2 frames.
type V2Key [32]byte

// NewV2Key allocates a V2Key.
func NewV2Key(in []byte) *V2Key {
	key := new(V2Key)
	copy(key[:], in)
	return key
}

// V2Signature is a V2 frame signature.
type V2Signature [6]byte

// V2Frame is a Mavlink V2 frame.
type V2Frame struct {
	IncompatibilityFlag byte
	CompatibilityFlag   byte
	SequenceID          byte
	SystemID            byte
	ComponentID         byte
	Message             message.Message
	Checksum            uint16
	SignatureLinkID     byte
	SignatureTimestamp  uint64
	Signature           *V2Signature
}

// GetSystemID implements Frame.
func (f V2Frame) GetSystemID() byte {
	return f.SystemID
}

// GetComponentID implements Frame.
func (f V2Frame) GetComponentID() byte {
	return f.ComponentID
}

// GetSequenceID implements Frame.
func (f V2Frame) GetSequenceID() byte {
	return f.SequenceID
}

// GetMessage implements Frame.
func (f V2Frame) GetMessage() message.Message {
	return f.Message
}

// IsSigned checks whether the frame contains a signature. It does not validate the signature.
func (f V2Frame) IsSigned() bool {
	return (f.IncompatibilityFlag & V2FlagSigned) != 0
}

// GetChecksum implements Frame.
func (f V2Frame) GetChecksum() uint16 {
	return f.Checksum
}

// GenerateChecksum implements Frame.
func (f V2Frame) GenerateChecksum(crcExtra byte) uint16 {
	msg := f.GetMessage().(*message.MessageRaw)
	h := x25.New()

	buf := make([]byte, 3)
	h.Write([]byte{byte(len(msg.Payload))})
	h.Write([]byte{f.IncompatibilityFlag})
	h.Write([]byte{f.CompatibilityFlag})
	h.Write([]byte{f.SequenceID})
	h.Write([]byte{f.SystemID})
	h.Write([]byte{f.ComponentID})
	uint24Encode(buf, msg.ID)
	h.Write(buf)
	h.Write(msg.Payload)

	h.Write([]byte{crcExtra})

	return h.Sum16()
}

// GenerateSignature generates the frame signature.
func (f V2Frame) GenerateSignature(key *V2Key) *V2Signature {
	msg := f.GetMessage().(*message.MessageRaw)
	h := sha256.New()

	// secret key
	h.Write(key[:])

	// the signature covers the whole message, excluding the signature itself
	buf := make([]byte, 6)
	h.Write([]byte{V2MagicByte})
	h.Write([]byte{byte(len(msg.Payload))})
	h.Write([]byte{f.IncompatibilityFlag})
	h.Write([]byte{f.CompatibilityFlag})
	h.Write([]byte{f.SequenceID})
	h.Write([]byte{f.SystemID})
	h.Write([]byte{f.ComponentID})
	uint24Encode(buf, msg.GetID())
	h.Write(buf[:3])
	h.Write(msg.Payload)
	binary.LittleEndian.PutUint16(buf, f.Checksum)
	h.Write(buf[:2])
	h.Write([]byte{f.SignatureLinkID})
	uint48Encode(buf, f.SignatureTimestamp)
	h.Write(buf)

	sig := new(V2Signature)
	copy(sig[:], h.Sum(nil)[:6])
	return sig
}

func (f *V2Frame) decode(br *bufio.Reader) error {
	// header
	buf, err := peekAndDiscard(br, 9)
	if err != nil {
		return err
	}
	msgLen := buf[0]
	f.IncompatibilityFlag = buf[1]
	f.CompatibilityFlag = buf[2]
	f.SequenceID = buf[3]
	f.SystemID = buf[4]
	f.ComponentID = buf[5]
	msgID := uint24Decode(buf[6:])

	// discard frame if incompatibility flag is not understood, as in recommendations
	if f.IncompatibilityFlag != 0 && f.IncompatibilityFlag != V2FlagSigned {
		return fmt.Errorf("unknown incompatibility flag: %d", f.IncompatibilityFlag)
	}

	// message
	var msgEncoded []byte
	if msgLen > 0 {
		msgEncoded = make([]byte, msgLen)
		_, err = io.ReadFull(br, msgEncoded)
		if err != nil {
			return err
		}
	}
	f.Message = &message.MessageRaw{
		ID:      msgID,
		Payload: msgEncoded,
	}

	// checksum
	buf, err = peekAndDiscard(br, 2)
	if err != nil {
		return err
	}
	f.Checksum = binary.LittleEndian.Uint16(buf)

	// signature
	if f.IsSigned() {
		buf, err := peekAndDiscard(br, 13)
		if err != nil {
			return err
		}
		f.SignatureLinkID = buf[0]
		f.SignatureTimestamp = uint48Decode(buf[1:])
		f.Signature = new(V2Signature)
		copy(f.Signature[:], buf[7:])
	}

	return nil
}

func (f V2Frame) encodeTo(buf []byte, msgEncoded []byte) (int, error) {
	msgLen := len(msgEncoded)

	// header
	buf[0] = V2MagicByte
	buf[1] = byte(msgLen)
	buf[2] = f.IncompatibilityFlag
	buf[3] = f.CompatibilityFlag
	buf[4] = f.SequenceID
	buf[5] = f.SystemID
	buf[6] = f.ComponentID
	uint24Encode(buf[7:], f.Message.GetID())
	n := 10

	// message
	if msgLen > 0 {
		n += copy(buf[n:], msgEncoded)
	}

	// checksum
	binary.LittleEndian.PutUint16(buf[n:], f.Checksum)
	n += 2

	// signature
	if f.IsSigned() {
		buf[n] = f.SignatureLinkID
		n++
		uint48Encode(buf[n:], f.SignatureTimestamp)
		n += 6
		n += copy(buf[n:], f.Signature[:])
	}

	return n, nil
}
