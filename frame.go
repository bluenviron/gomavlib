package gomavlib

const (
	v1MagicByte = '\xFE'
	v2MagicByte = '\xFD'
	flagSigned  = 0x01
)

type Frame interface {
	GetVersion() int
	GetSystemId() byte
	GetComponentId() byte
	GetMessage() Message
	GetChecksum() uint16
}

type FrameV1 struct {
	SequenceId  byte
	SystemId    byte
	ComponentId byte
	Message     Message
	Checksum    uint16
}

func (f *FrameV1) GetVersion() int {
	return 1
}

func (f *FrameV1) GetSystemId() byte {
	return f.SystemId
}

func (f *FrameV1) GetComponentId() byte {
	return f.ComponentId
}

func (f *FrameV1) GetMessage() Message {
	return f.Message
}

func (f *FrameV1) GetChecksum() uint16 {
	return f.Checksum
}

type FrameV2 struct {
	IncompatibilityFlag byte
	CompatibilityFlag   byte
	SequenceId          byte
	SystemId            byte
	ComponentId         byte
	Message             Message
	Checksum            uint16
	SignatureLinkId     byte
	SignatureTimestamp  uint64
	Signature           Signature
}

func (f *FrameV2) GetVersion() int {
	return 2
}

func (f *FrameV2) GetSystemId() byte {
	return f.SystemId
}

func (f *FrameV2) GetComponentId() byte {
	return f.ComponentId
}

func (f *FrameV2) GetMessage() Message {
	return f.Message
}

func (f *FrameV2) GetChecksum() uint16 {
	return f.Checksum
}

func (f *FrameV2) isSigned() bool {
	return (f.IncompatibilityFlag & flagSigned) != 0
}

type SignatureKey [32]byte

func NewSignatureKey(in []byte) *SignatureKey {
	var key SignatureKey
	copy(key[:], in)
	return &key
}

type Signature [6]byte
