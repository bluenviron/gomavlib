package frame

import (
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/message"
)

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
		return nil, fmt.Errorf("OutSystemID must be greater than one")
	}
	if conf.OutComponentID < 1 {
		conf.OutComponentID = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	return &Writer{
		conf: conf,
		bw:   make([]byte, bufferSize),
	}, nil
}

// WriteMessage writes a Message.
// The Message is wrapped into a Frame whose fields are filled automatically.
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

	// fill SequenceID, SystemID, ComponentID, CompatibilityFlag, IncompatibilityFlag
	switch ff := fr.(type) {
	case *V1Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID

	case *V2Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID

		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0
		if w.conf.OutKey != nil {
			ff.IncompatibilityFlag |= V2FlagSigned
		}
	}

	w.curWriteSequenceID++

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if w.conf.DialectRW == nil {
			return fmt.Errorf("dialect is nil")
		}

		mp := w.conf.DialectRW.GetMessage(fr.GetMessage().GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		_, isV2 := fr.(*V2Frame)
		byts := mp.Write(fr.GetMessage(), isV2)

		msgRaw := &message.MessageRaw{
			ID:      fr.GetMessage().GetID(),
			Payload: byts,
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
			ff.Checksum = ff.generateChecksum(mp.CRCExtra())
		case *V2Frame:
			ff.Checksum = ff.generateChecksum(mp.CRCExtra())
		}
	}

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if ff, ok := fr.(*V2Frame); ok && w.conf.OutKey != nil {
		ff.SignatureLinkID = w.conf.OutSignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.genSignature(w.conf.OutKey)
	}

	return w.writeFrameInner(fr, fr.GetMessage().(*message.MessageRaw))
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
			return fmt.Errorf("dialect is nil")
		}

		mp := w.conf.DialectRW.GetMessage(m.GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		_, isV2 := fr.(*V2Frame)
		byts := mp.Write(m, isV2)

		// do not touch Message
		// in such way that the frame can be encoded by other parsers in parallel
		m = &message.MessageRaw{
			ID:      m.GetID(),
			Payload: byts,
		}
	}

	return w.writeFrameInner(fr, m.(*message.MessageRaw))
}

func (w *Writer) writeFrameInner(fr Frame, m *message.MessageRaw) error {
	n, err := fr.encodeTo(w.bw, m.Payload)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.conf.Writer.Write(w.bw[:n])
	return err
}
