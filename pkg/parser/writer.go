package parser

import (
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
)

// WriterConf is the configuration of a Writer.
type WriterConf struct {
	// the underlying bytes writer.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	DialectDE *dialect.DecEncoder

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
	OutKey *frame.V2Key
}

// Writer is a Mavlink writer.
type Writer struct {
	conf               WriterConf
	writeBuffer        []byte
	curWriteSequenceID byte
}

// NewWriter allocates a writer.
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
		conf:        conf,
		writeBuffer: make([]byte, 0, bufferSize),
	}, nil
}

// WriteMessage writes a Message.
// It must not be called by multiple routines in parallel.
func (w *Writer) WriteMessage(m msg.Message) error {
	if w.conf.OutVersion == V1 {
		return w.writeFrameAndFill(&frame.V1Frame{Message: m})
	}
	return w.writeFrameAndFill(&frame.V2Frame{Message: m})
}

func (w *Writer) writeFrameAndFill(fr frame.Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// do not touch the original frame, but work with a separate object
	// in such way that the frame can be encoded by other parsers in parallel
	safeFrame := fr.Clone()

	// fill SequenceID, SystemID, ComponentID
	switch ff := safeFrame.(type) {
	case *frame.V1Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID
	case *frame.V2Frame:
		ff.SequenceID = w.curWriteSequenceID
		ff.SystemID = w.conf.OutSystemID
		ff.ComponentID = w.conf.OutComponentID
	}
	w.curWriteSequenceID++

	// fill CompatibilityFlag, IncompatibilityFlag if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok {
		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0

		if w.conf.OutKey != nil {
			ff.IncompatibilityFlag |= frame.V2FlagSigned
		}
	}

	// encode message if it is not already encoded
	if _, ok := safeFrame.GetMessage().(*msg.MessageRaw); !ok {
		if w.conf.DialectDE == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := w.conf.DialectDE.MessageDEs[safeFrame.GetMessage().GetID()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := safeFrame.(*frame.V2Frame)
		byt, err := mp.Encode(safeFrame.GetMessage(), isV2)
		if err != nil {
			return err
		}

		msgRaw := &msg.MessageRaw{safeFrame.GetMessage().GetID(), byt} //nolint:govet
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Message = msgRaw
		case *frame.V2Frame:
			ff.Message = msgRaw
		}

		// fill checksum
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Checksum = ff.GenChecksum(w.conf.DialectDE.MessageDEs[ff.GetMessage().GetID()].CRCExtra())
		case *frame.V2Frame:
			ff.Checksum = ff.GenChecksum(w.conf.DialectDE.MessageDEs[ff.GetMessage().GetID()].CRCExtra())
		}
	}

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok && w.conf.OutKey != nil {
		ff.SignatureLinkID = w.conf.OutSignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.GenSignature(w.conf.OutKey)
	}

	return w.WriteFrame(safeFrame)
}

// WriteFrame writes a Frame.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (w *Writer) WriteFrame(fr frame.Frame) error {
	m := fr.GetMessage()
	if m == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := m.(*msg.MessageRaw); !ok {
		if w.conf.DialectDE == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := w.conf.DialectDE.MessageDEs[m.GetID()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := fr.(*frame.V2Frame)
		byt, err := mp.Encode(m, isV2)
		if err != nil {
			return err
		}

		// do not touch frame.Message
		// in such way that the frame can be encoded by other parsers in parallel
		m = &msg.MessageRaw{m.GetID(), byt} //nolint:govet
	}

	buf, err := fr.Encode(w.writeBuffer, m.(*msg.MessageRaw).Content)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.conf.Writer.Write(buf)
	return err
}
