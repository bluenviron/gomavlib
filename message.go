package gomavlib

type Message interface {
	GetId() uint32
}

type MessageRaw struct {
	Id      uint32
	Content []byte
}

func (m *MessageRaw) GetId() uint32 {
	return m.Id
}
