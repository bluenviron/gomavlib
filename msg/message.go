// Package msg contains the Message definition and utilities to encode and
// decode messages.
package msg

// Message is the interface that must be implemented by all Mavlink messages.
// Furthermore, any message must be labeled "MessageNameOfMessage".
type Message interface {
	GetId() uint32
}

// MessageRaw is a special struct that contains a byte-encoded message,
// available in Content. It is used:
//
// * as intermediate step in the encoding/decoding process
//
// * when the parser receives an unknown message
//
type MessageRaw struct {
	Id      uint32
	Content []byte
}

// GetId implements the Message interface.
func (m *MessageRaw) GetId() uint32 {
	return m.Id
}
