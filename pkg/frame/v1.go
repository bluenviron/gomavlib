package frame

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/aler9/gomavlib/pkg/msg"
	"github.com/aler9/gomavlib/pkg/x25"
)

const (
	// V1MagicByte is the magic byte of a V1 frame.
	V1MagicByte = 0xFE
)

// V1Frame is a Mavlink V1 frame.
type V1Frame struct {
	SequenceId  byte
	SystemId    byte
	ComponentId byte
	Message     msg.Message
	Checksum    uint16
}

// Clone implements the Frame interface.
func (f *V1Frame) Clone() Frame {
	return &V1Frame{
		SequenceId:  f.SequenceId,
		SystemId:    f.SystemId,
		ComponentId: f.ComponentId,
		Message:     f.Message,
		Checksum:    f.Checksum,
	}
}

// GetSystemId implements the Frame interface.
func (f *V1Frame) GetSystemId() byte {
	return f.SystemId
}

// GetComponentId implements the Frame interface.
func (f *V1Frame) GetComponentId() byte {
	return f.ComponentId
}

// GetMessage implements the Frame interface.
func (f *V1Frame) GetMessage() msg.Message {
	return f.Message
}

// GetChecksum implements the Frame interface.
func (f *V1Frame) GetChecksum() uint16 {
	return f.Checksum
}

// Decode implements the Frame interface.
func (f *V1Frame) Decode(br *bufio.Reader) error {
	// header
	buf, err := br.Peek(5)
	if err != nil {
		return err
	}
	br.Discard(5)
	msgLen := buf[0]
	f.SequenceId = buf[1]
	f.SystemId = buf[2]
	f.ComponentId = buf[3]
	msgId := buf[4]

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
		Id:      uint32(msgId),
		Content: msgEncoded,
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

// Encode implements the Frame interface.
func (f *V1Frame) Encode(buf []byte, msgEncoded []byte) ([]byte, error) {
	if f.Message.GetId() > 0xFF {
		return nil, fmt.Errorf("cannot send a message with an id > 0xFF inside a V1 frame")
	}

	msgLen := len(msgEncoded)
	bufLen := 6 + msgLen + 2
	buf = buf[:bufLen]

	// header
	buf[0] = V1MagicByte
	buf[1] = byte(msgLen)
	buf[2] = f.SequenceId
	buf[3] = f.SystemId
	buf[4] = f.ComponentId
	buf[5] = byte(f.Message.GetId())

	// message
	if msgLen > 0 {
		copy(buf[6:], msgEncoded)
	}

	// checksum
	binary.LittleEndian.PutUint16(buf[6+msgLen:], f.Checksum)

	return buf, nil
}

// GenChecksum implements the Frame interface.
func (f *V1Frame) GenChecksum(crcExtra byte) uint16 {
	msg := f.GetMessage().(*msg.MessageRaw)
	h := x25.New()

	h.Write([]byte{byte(len(msg.Content))})
	h.Write([]byte{f.SequenceId})
	h.Write([]byte{f.SystemId})
	h.Write([]byte{f.ComponentId})
	h.Write([]byte{byte(msg.Id)})
	h.Write(msg.Content)

	h.Write([]byte{crcExtra})

	return h.Sum16()
}
