// Package msg contains the Message definition and utilities to encode and
// decode messages.
package msg

// Message is the interface that must be implemented by all Mavlink messages.
// Furthermore, any message must be labeled "MessageNameOfMessage".
type Message interface {
	GetID() uint32
}

// MessageRaw is a special struct that contains an unencoded message.
// It is used:
//
// * as intermediate step in the encoding/decoding process
//
// * when the parser receives an unknown message
//
type MessageRaw struct {
	ID      uint32
	Content []byte
}

// GetID implements the Message interface.
func (m *MessageRaw) GetID() uint32 {
	return m.ID
}
