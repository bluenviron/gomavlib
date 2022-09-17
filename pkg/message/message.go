// Package message contains the message definition and its parser.
package message

// Message is the interface that must be implemented by all Mavlink messages.
// Furthermore, any message struct must be labeled "MessageNameOfMessage".
type Message interface {
	GetID() uint32
}

// MessageRaw is a special struct that contains an unencoded message.
// It is used:
//
// * as intermediate step in the encoding/decoding process
//
// * when the parser receives an unknown message
type MessageRaw struct { //nolint:revive
	ID      uint32
	Payload []byte
}

// GetID implements the Message interface.
func (m *MessageRaw) GetID() uint32 {
	return m.ID
}
