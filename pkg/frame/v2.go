package frame

import (
	"bufio"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/aler9/gomavlib/pkg/x25"
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
	Message             msg.Message
	Checksum            uint16
	SignatureLinkID     byte
	SignatureTimestamp  uint64
	Signature           *V2Signature
}

// Clone implements the Frame interface.
func (f *V2Frame) Clone() Frame {
	return &V2Frame{
		IncompatibilityFlag: f.IncompatibilityFlag,
		CompatibilityFlag:   f.CompatibilityFlag,
		SequenceID:          f.SequenceID,
		SystemID:            f.SystemID,
		ComponentID:         f.ComponentID,
		Message:             f.Message,
		Checksum:            f.Checksum,
		SignatureLinkID:     f.SignatureLinkID,
		SignatureTimestamp:  f.SignatureTimestamp,
		Signature:           f.Signature,
	}
}

// GetSystemID implements the Frame interface.
func (f *V2Frame) GetSystemID() byte {
	return f.SystemID
}

// GetComponentID implements the Frame interface.
func (f *V2Frame) GetComponentID() byte {
	return f.ComponentID
}

// GetMessage implements the Frame interface.
func (f *V2Frame) GetMessage() msg.Message {
	return f.Message
}

// GetChecksum implements the Frame interface.
func (f *V2Frame) GetChecksum() uint16 {
	return f.Checksum
}

// IsSigned checks whether the frame contains a signature. It does not validate the signature.
func (f *V2Frame) IsSigned() bool {
	return (f.IncompatibilityFlag & V2FlagSigned) != 0
}

// Decode implements the Frame interface.
func (f *V2Frame) Decode(br *bufio.Reader) error {
	// header
	buf, err := br.Peek(9)
	if err != nil {
		return err
	}
	br.Discard(9)
	msgLen := buf[0]
	f.IncompatibilityFlag = buf[1]
	f.CompatibilityFlag = buf[2]
	f.SequenceID = buf[3]
	f.SystemID = buf[4]
	f.ComponentID = buf[5]
	msgID := uint24Decode(buf[6:])

	// discard frame if incompatibility flag is not understood, as in recommendations
	if f.IncompatibilityFlag != 0 && f.IncompatibilityFlag != V2FlagSigned {
		return fmt.Errorf("unknown incompatibility flag (%d)", f.IncompatibilityFlag)
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
	f.Message = &msg.MessageRaw{
		ID:      msgID,
		Content: msgEncoded,
	}

	// checksum
	buf, err = br.Peek(2)
	if err != nil {
		return err
	}
	br.Discard(2)
	f.Checksum = binary.LittleEndian.Uint16(buf)

	// signature
	if f.IsSigned() {
		buf, err := br.Peek(13)
		if err != nil {
			return err
		}
		br.Discard(13)
		f.SignatureLinkID = buf[0]
		f.SignatureTimestamp = uint48Decode(buf[1:])
		f.Signature = new(V2Signature)
		copy(f.Signature[:], buf[7:])
	}

	return nil
}

// Encode implements the Frame interface.
func (f *V2Frame) Encode(buf []byte, msgEncoded []byte) ([]byte, error) {
	msgLen := len(msgEncoded)
	bufLen := 10 + msgLen + 2
	if f.IsSigned() {
		bufLen += 13
	}
	buf = buf[:bufLen]

	// header
	buf[0] = V2MagicByte
	buf[1] = byte(msgLen)
	buf[2] = f.IncompatibilityFlag
	buf[3] = f.CompatibilityFlag
	buf[4] = f.SequenceID
	buf[5] = f.SystemID
	buf[6] = f.ComponentID
	uint24Encode(buf[7:], f.Message.GetID())

	// message
	if msgLen > 0 {
		copy(buf[10:], msgEncoded)
	}

	// checksum
	binary.LittleEndian.PutUint16(buf[10+msgLen:], f.Checksum)

	// signature
	if f.IsSigned() {
		buf[12+msgLen] = f.SignatureLinkID
		uint48Encode(buf[13+msgLen:], f.SignatureTimestamp)
		copy(buf[19+msgLen:], f.Signature[:])
	}

	return buf, nil
}

// GenChecksum implements the Frame interface.
func (f *V2Frame) GenChecksum(crcExtra byte) uint16 {
	msg := f.GetMessage().(*msg.MessageRaw)
	h := x25.New()

	buf := make([]byte, 3)
	h.Write([]byte{byte(len(msg.Content))})
	h.Write([]byte{f.IncompatibilityFlag})
	h.Write([]byte{f.CompatibilityFlag})
	h.Write([]byte{f.SequenceID})
	h.Write([]byte{f.SystemID})
	h.Write([]byte{f.ComponentID})
	h.Write(uint24Encode(buf, msg.ID))
	h.Write(msg.Content)

	h.Write([]byte{crcExtra})

	return h.Sum16()
}

// GenSignature generates a signature with the given key.
func (f *V2Frame) GenSignature(key *V2Key) *V2Signature {
	msg := f.GetMessage().(*msg.MessageRaw)
	h := sha256.New()

	// secret key
	h.Write(key[:])

	// the signature covers the whole message, excluding the signature itself
	buf := make([]byte, 6)
	h.Write([]byte{V2MagicByte})
	h.Write([]byte{byte(len(msg.Content))})
	h.Write([]byte{f.IncompatibilityFlag})
	h.Write([]byte{f.CompatibilityFlag})
	h.Write([]byte{f.SequenceID})
	h.Write([]byte{f.SystemID})
	h.Write([]byte{f.ComponentID})
	h.Write(uint24Encode(buf, f.Message.GetID()))
	h.Write(msg.Content)
	binary.LittleEndian.PutUint16(buf, f.Checksum)
	h.Write(buf[:2])
	h.Write([]byte{f.SignatureLinkID})
	h.Write(uint48Encode(buf, f.SignatureTimestamp))

	sig := new(V2Signature)
	copy(sig[:], h.Sum(nil)[:6])
	return sig
}
