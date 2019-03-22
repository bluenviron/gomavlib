package gomavlib

// Message is the interface that all mavlink messages must implements.
// Furthermore, message structs must be labeled MessageNameOfMessage.
type Message interface {
	GetId() uint32
}

// MessageRaw is a special struct that contains a byte-encoded message,
// available in Content. Is is used:
// * as intermediate step in the encoding/decoding process
// * when the parser receives an unknown message
type MessageRaw struct {
	Id      uint32
	Content []byte
}

// GetId implements the message interface.
func (m *MessageRaw) GetId() uint32 {
	return m.Id
}
