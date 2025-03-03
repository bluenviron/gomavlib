// Package streamwriter contains a message stream writer.
package streamwriter

import (
	"fmt"
	"time"

	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 0o1, 0o1, 0, 0, 0, 0, time.UTC)

func encodeMessageInFrame(fr frame.Frame, mp *message.ReadWriter) {
	_, isV2 := fr.(*frame.V2Frame)
	msgRaw := mp.Write(fr.GetMessage(), isV2)

	switch ff := fr.(type) {
	case *frame.V1Frame:
		ff.Message = msgRaw
	case *frame.V2Frame:
		ff.Message = msgRaw
	}
}

// Version is a Mavlink version.
type Version int

const (
	// V1 is Mavlink 1.0
	V1 Version = 1

	// V2 is Mavlink 2.0
	V2 Version = 2
)

// String implements fmt.Stringer.
func (v Version) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}

// Writer is a message stream writer.
// It allows to send out messages, wrapped in frames.
// Frame fields are filled automatically.
// It must not be called by multiple routines in parallel.
type Writer struct {
	// frame writer.
	FrameWriter *frame.Writer

	// Mavlink version used to encode outgoing frames.
	Version Version
	// system id, added to every outgoing frame and used to identify this
	// node in the network.
	SystemID byte
	// (optional) component id, added to every outgoing frame, defaults to 1.
	ComponentID byte
	// (optional) value to insert into the signature link id.
	// This feature requires v2 frames.
	SignatureLinkID byte
	// (optional) secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	Key *frame.V2Key

	//
	// private
	//

	nextSeqNumber byte
}

// Initialize initializes a Writer.
func (w *Writer) Initialize() error {
	if w.Version == 0 {
		return fmt.Errorf("OutVersion not provided")
	}
	if w.SystemID < 1 {
		return fmt.Errorf("OutSystemID must be greater than one")
	}
	if w.ComponentID < 1 {
		w.ComponentID = 1
	}
	if w.Key != nil && w.Version != V2 {
		return fmt.Errorf("OutKey requires V2 frames")
	}

	return nil
}

// Write writes a message.
func (w *Writer) Write(msg message.Message) error {
	if w.Version == V1 {
		return w.writeInner(&frame.V1Frame{Message: msg})
	}
	return w.writeInner(&frame.V2Frame{Message: msg})
}

func (w *Writer) writeInner(fr frame.Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// fill SequenceNumber, SystemID, ComponentID, CompatibilityFlag, IncompatibilityFlag
	switch ff := fr.(type) {
	case *frame.V1Frame:
		ff.SequenceNumber = w.nextSeqNumber
		ff.SystemID = w.SystemID
		ff.ComponentID = w.ComponentID

	case *frame.V2Frame:
		ff.SequenceNumber = w.nextSeqNumber
		ff.SystemID = w.SystemID
		ff.ComponentID = w.ComponentID

		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0
		if w.Key != nil {
			ff.IncompatibilityFlag |= frame.V2FlagSigned
		}
	}

	w.nextSeqNumber++

	if w.FrameWriter.DialectRW == nil {
		return fmt.Errorf("dialect is nil")
	}

	mp := w.FrameWriter.DialectRW.GetMessage(fr.GetMessage().GetID())
	if mp == nil {
		return fmt.Errorf("message is not in the dialect")
	}

	// encode message if it is not already encoded
	if _, ok := fr.GetMessage().(*message.MessageRaw); !ok {
		encodeMessageInFrame(fr, mp)
	}

	// fill checksum
	switch ff := fr.(type) {
	case *frame.V1Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	case *frame.V2Frame:
		ff.Checksum = ff.GenerateChecksum(mp.CRCExtra())
	}

	// fill SignatureLinkID, SignatureTimestamp, Signature if v2
	if ff, ok := fr.(*frame.V2Frame); ok && w.Key != nil {
		ff.SignatureLinkID = w.SignatureLinkID
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.GenerateSignature(w.Key)
	}

	return w.FrameWriter.Write(fr)
}
