package frame

import (
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/message"
)

// WriterOutVersion is a Mavlink version.
type WriterOutVersion int

const (
	// V1 is Mavlink 1.0
	V1 WriterOutVersion = 1

	// V2 is Mavlink 2.0
	V2 WriterOutVersion = 2
)

// String implements fmt.Stringer.
func (v WriterOutVersion) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}

// WriterConf is the configuration of a Writer.
type WriterConf struct {
	// the underlying bytes writer.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key
}

// Writer is a Frame writer.
type Writer struct {
	conf               WriterConf
	bw                 []byte
	curWriteSequenceID byte
}

// NewWriter allocates a Writer.
func NewWriter(conf WriterConf) (*Writer, error) {
	if conf.Writer == nil {
		return nil, fmt.Errorf("Writer not provided")
	}

	if conf.OutVersion == 0 {
		return nil, fmt.Errorf("OutVersion not provided")
	}
	if conf.OutSystemID < 1 {
		return nil, fmt.Errorf("SystemID must be >= 1")
	}
	if conf.OutComponentID < 1 {
		conf.OutComponentID = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	return &Writer{
		conf: conf,
		bw:   make([]byte, 0, bufferSize),
	}, nil
}

// WriteMessage writes a Message.
// It must not be called by multiple routines in parallel.
func (w *Writer) WriteMessage(m message.Message) error {
	if w.conf.OutVersion == V1 {
		return w.writeFrameAndFill(&V1Frame{Message: m})
	}
	return w.writeFrameAndFill(&V2Frame{Message: m})
}

func (w *Writer) writeFrameAndFill(fr Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// do not touch the original frame, but work with a separate object
	// in such way that the frame can be encoded by other parsers in parallel
	fr = fr.Clone()

	// fill SequenceID, SystemID, ComponentID
	switch ff := fr.(type) {
	case *V1Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID
	case *V2Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID
	}
	w.curWriteSequenceID++

	// fill CompatibilityFlag, IncompatibilityFlag if v2
	if ff, ok := fr.(*V2Frame); ok {
		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0

		if w.conf.OutKey != nil {
			ff.IncompatibilityFlag |= V2FlagSigned
		}
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if w.conf.DialectRW == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := w.conf.DialectRW.MessageDEs[fr.GetMessage().GetID()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := fr.(*V2Frame)
		byt, err := mp.Write(fr.GetMessage(), isV2)
		if err != nil {
			return err
		}

		msgRaw := &message.MessageRaw{
			ID:      fr.GetMessage().GetID(),
			Payload: byt,
		}
		switch ff := fr.(type) {
		case *V1Frame:
			ff.Message = msgRaw
		case *V2Frame:
			ff.Message = msgRaw
		}

		// fill checksum
		switch ff := fr.(type) {
		case *V1Frame:
			ff.Checksum = ff.genChecksum(w.conf.DialectRW.MessageDEs[ff.GetMessage().GetID()].CRCExtra())
		case *V2Frame:
			ff.Checksum = ff.genChecksum(w.conf.DialectRW.MessageDEs[ff.GetMessage().GetID()].CRCExtra())
		}
	}

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if ff, ok := fr.(*V2Frame); ok && w.conf.OutKey != nil {
		ff.SignatureLinkID = w.conf.OutSignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.genSignature(w.conf.OutKey)
	}

	return w.WriteFrame(fr)
}

// WriteFrame writes a Frame.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (w *Writer) WriteFrame(fr Frame) error {
	m := fr.GetMessage()
	if m == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := m.(*message.MessageRaw); !ok {
		if w.conf.DialectRW == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := w.conf.DialectRW.MessageDEs[m.GetID()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := fr.(*V2Frame)
		byt, err := mp.Write(m, isV2)
		if err != nil {
			return err
		}

		// do not touch Message
		// in such way that the frame can be encoded by other parsers in parallel
		m = &message.MessageRaw{m.GetID(), byt} //nolint:govet
	}

	buf, err := fr.encode(w.bw, m.(*message.MessageRaw).Payload)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.conf.Writer.Write(buf)
	return err
}
