package frame

import (
	"fmt"
	"io"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// WriterConf is the configuration of a Writer.
//
// Deprecated: configuration has been moved inside Writer.
type WriterConf struct {
	// underlying bytes writer.
	Writer io.Writer

	// (optional) dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key
}

// NewWriter allocates a Writer.
//
// Deprecated: replaced by Writer.Initialize().
func NewWriter(conf WriterConf) (*Writer, error) {
	w := &Writer{
		ByteWriter:         conf.Writer,
		DialectRW:          conf.DialectRW,
		OutVersion:         conf.OutVersion,
		OutSystemID:        conf.OutSystemID,
		OutComponentID:     conf.OutComponentID,
		OutSignatureLinkID: conf.OutSignatureLinkID,
		OutKey:             conf.OutKey,
	}
	err := w.Initialize()
	return w, err
}

// Writer is a Frame writer.
type Writer struct {
	// underlying byte writer.
	ByteWriter io.Writer

	// (optional) dialect which contains the messages that will be written.
	DialectRW *dialect.ReadWriter

	// Mavlink version used to encode messages.
	OutVersion WriterOutVersion
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	OutComponentID byte
	// (optional) value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *V2Key

	//
	// private
	//

	bw            []byte
	nextSeqNumber byte
}

// Initialize allocates a Writer.
func (w *Writer) Initialize() error {
	if w.ByteWriter == nil {
		return fmt.Errorf("ByteWriter not provided")
	}

	if w.OutVersion == 0 {
		return fmt.Errorf("OutVersion not provided")
	}
	if w.OutSystemID < 1 {
		return fmt.Errorf("OutSystemID must be greater than one")
	}
	if w.OutComponentID < 1 {
		w.OutComponentID = 1
	}
	if w.OutKey != nil && w.OutVersion != V2 {
		return fmt.Errorf("OutKey requires V2 frames")
	}

	w.bw = make([]byte, bufferSize)

	return nil
}

// WriteMessage writes a Message.
// The Message is wrapped into a Frame whose fields are filled automatically.
// It must not be called by multiple routines in parallel.
func (w *Writer) WriteMessage(m message.Message) error {
	if w.OutVersion == V1 {
		return w.writeFrameAndFill(&V1Frame{Message: m})
	}
	return w.writeFrameAndFill(&V2Frame{Message: m})
}

func (w *Writer) writeFrameAndFill(fr Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// fill SequenceNumber, SystemID, ComponentID, CompatibilityFlag, IncompatibilityFlag
	switch ff := fr.(type) {
	case *V1Frame:
		ff.SequenceNumber = w.nextSeqNumber
		ff.SystemID = w.OutSystemID
		ff.ComponentID = w.OutComponentID

	case *V2Frame:
		ff.SequenceNumber = w.nextSeqNumber
		ff.SystemID = w.OutSystemID
		ff.ComponentID = w.OutComponentID

		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0
		if w.OutKey != nil {
			ff.IncompatibilityFlag |= V2FlagSigned
		}
	}

	w.nextSeqNumber++

	if w.DialectRW == nil {
		return fmt.Errorf("dialect is nil")
	}

	mp := w.DialectRW.GetMessage(fr.GetMessage().GetID())
	if mp == nil {
		return fmt.Errorf("message is not in the dialect")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		w.encodeMessageInFrame(fr, mp)
	}

	// fill checksum
	switch ff := fr.(type) {
	case *V1Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	case *V2Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	}

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if ff, ok := fr.(*V2Frame); ok && w.OutKey != nil {
		ff.SignatureLinkID = w.OutSignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.GenerateSignature(w.OutKey)
	}

	return w.writeFrameInner(fr)
}

// WriteFrame writes a Frame.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (w *Writer) WriteFrame(fr Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		if w.DialectRW == nil {
			return fmt.Errorf("dialect is nil")
		}

		mp := w.DialectRW.GetMessage(fr.GetMessage().GetID())
		if mp == nil {
			return fmt.Errorf("message is not in the dialect")
		}

		w.encodeMessageInFrame(fr, mp)
	}

	return w.writeFrameInner(fr)
}

func (w *Writer) encodeMessageInFrame(fr Frame, mp *message.ReadWriter) {
	_, isV2 := fr.(*V2Frame)
	msgRaw := mp.Write(fr.GetMessage(), isV2)

	switch ff := fr.(type) {
	case *V1Frame:
		ff.Message = msgRaw
	case *V2Frame:
		ff.Message = msgRaw
	}
}

func (w *Writer) writeFrameInner(fr Frame) error {
	n, err := fr.marshalTo(w.bw, fr.GetMessage().(*message.MessageRaw).Payload)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = w.ByteWriter.Write(w.bw[:n])
	return err
}
