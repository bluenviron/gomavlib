package gomavlib

const (
	v1MagicByte = 0xfe
	v2MagicByte = 0xfd
	flagSigned  = 0x01
)

// Frame is the interface implemented by frames of every supported version.
type Frame interface {
	// the frame version.
	GetVersion() int
	// the system id of the author of the frame.
	GetSystemId() byte
	// the component id of the author of the frame.
	GetComponentId() byte
	// the message encapsuled in the frame.
	GetMessage() Message
	// the frame checksum.
	GetChecksum() uint16
	// generate a clone of the frame
	Clone() Frame
}

// FrameV1 represents a 1.0 frame.
type FrameV1 struct {
	SequenceId  byte
	SystemId    byte
	ComponentId byte
	Message     Message
	Checksum    uint16
}

// Clone is part of the Frame interface.
func (f *FrameV1) Clone() Frame {
	return &FrameV1{
		SequenceId:  f.SequenceId,
		SystemId:    f.SystemId,
		ComponentId: f.ComponentId,
		Message:     f.Message,
		Checksum:    f.Checksum,
	}
}

// GetVersion is part of the Frame interface.
func (f *FrameV1) GetVersion() int {
	return 1
}

// GetSystemId is part of the Frame interface.
func (f *FrameV1) GetSystemId() byte {
	return f.SystemId
}

// GetComponentId is part of the Frame interface.
func (f *FrameV1) GetComponentId() byte {
	return f.ComponentId
}

// GetMessage is part of the Frame interface.
func (f *FrameV1) GetMessage() Message {
	return f.Message
}

// GetChecksum is part of the Frame interface.
func (f *FrameV1) GetChecksum() uint16 {
	return f.Checksum
}

// FrameV2 represents a 2.0 frame.
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
	Signature           *Signature
}

// Clone is part of the Frame interface.
func (f *FrameV2) Clone() Frame {
	return &FrameV2{
		IncompatibilityFlag: f.IncompatibilityFlag,
		CompatibilityFlag:   f.CompatibilityFlag,
		SequenceId:          f.SequenceId,
		SystemId:            f.SystemId,
		ComponentId:         f.ComponentId,
		Message:             f.Message,
		Checksum:            f.Checksum,
		SignatureLinkId:     f.SignatureLinkId,
		SignatureTimestamp:  f.SignatureTimestamp,
		Signature:           f.Signature,
	}
}

// GetVersion is part of the Frame interface.
func (f *FrameV2) GetVersion() int {
	return 2
}

// GetSystemId is part of the Frame interface.
func (f *FrameV2) GetSystemId() byte {
	return f.SystemId
}

// GetComponentId is part of the Frame interface.
func (f *FrameV2) GetComponentId() byte {
	return f.ComponentId
}

// GetMessage is part of the Frame interface.
func (f *FrameV2) GetMessage() Message {
	return f.Message
}

// GetChecksum is part of the Frame interface.
func (f *FrameV2) GetChecksum() uint16 {
	return f.Checksum
}

// IsSigned checks whether the frame contains a signature. It does not validate the signature.
func (f *FrameV2) IsSigned() bool {
	return (f.IncompatibilityFlag & flagSigned) != 0
}

// SignatureKey is a key able to sign and validate V2 frames.
type SignatureKey [32]byte

// NewSignatureKey allocates a SignatureKey
func NewSignatureKey(in []byte) *SignatureKey {
	key := new(SignatureKey)
	copy(key[:], in)
	return key
}

// Signature is a V2 frame signature.
type Signature [6]byte
