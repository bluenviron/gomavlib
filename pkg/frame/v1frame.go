package frame

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/bluenviron/gomavlib/v2/pkg/message"
	"github.com/bluenviron/gomavlib/v2/pkg/x25"
)

const (
	// V1MagicByte is the magic byte of a V1 frame.
	V1MagicByte = 0xFE
)

// V1Frame is a Mavlink V1 frame.
type V1Frame struct {
	SequenceID  byte
	SystemID    byte
	ComponentID byte
	Message     message.Message
	Checksum    uint16
}

// GetSystemID implements Frame.
func (f V1Frame) GetSystemID() byte {
	return f.SystemID
}

// GetComponentID implements Frame.
func (f V1Frame) GetComponentID() byte {
	return f.ComponentID
}

// GetMessage implements Frame.
func (f V1Frame) GetMessage() message.Message {
	return f.Message
}

// GetChecksum implements Frame.
func (f V1Frame) GetChecksum() uint16 {
	return f.Checksum
}

// GenerateChecksum implements Frame.
func (f V1Frame) GenerateChecksum(crcExtra byte) uint16 {
	msg := f.GetMessage().(*message.MessageRaw)
	h := x25.New()

	h.Write([]byte{byte(len(msg.Payload))})
	h.Write([]byte{f.SequenceID})
	h.Write([]byte{f.SystemID})
	h.Write([]byte{f.ComponentID})
	h.Write([]byte{byte(msg.ID)})
	h.Write(msg.Payload)

	h.Write([]byte{crcExtra})

	return h.Sum16()
}

func (f *V1Frame) decode(br *bufio.Reader) error {
	// header
	buf, err := br.Peek(5)
	if err != nil {
		return err
	}
	br.Discard(5)
	msgLen := buf[0]
	f.SequenceID = buf[1]
	f.SystemID = buf[2]
	f.ComponentID = buf[3]
	msgID := buf[4]

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
		ID:      uint32(msgID),
		Payload: msgEncoded,
	}

	// checksum
	buf, err = br.Peek(2)
	if err != nil {
		return err
	}
	br.Discard(2)
	f.Checksum = binary.LittleEndian.Uint16(buf)

	return nil
}

func (f V1Frame) encodeTo(buf []byte, msgEncoded []byte) (int, error) {
	if f.Message.GetID() > 0xFF {
		return 0, fmt.Errorf("cannot send a message with an ID greater than 255 with a V1 frame")
	}

	msgLen := len(msgEncoded)

	// header
	buf[0] = V1MagicByte
	buf[1] = byte(msgLen)
	buf[2] = f.SequenceID
	buf[3] = f.SystemID
	buf[4] = f.ComponentID
	buf[5] = byte(f.Message.GetID())
	n := 6

	// message
	if msgLen > 0 {
		n += copy(buf[n:], msgEncoded)
	}

	// checksum
	binary.LittleEndian.PutUint16(buf[n:], f.Checksum)
	n += 2

	return n, nil
}
