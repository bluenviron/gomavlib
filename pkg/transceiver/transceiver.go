// Package transceiver implements a Mavlink transceiver.
package transceiver

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
)

const (
	bufferSize = 512 // frames cannot go beyond len(header) + 255 + len(check) + len(sig)
)

// 1st January 2015 GMT
var signatureReferenceDate = time.Date(2015, 01, 01, 0, 0, 0, 0, time.UTC)

// TransceiverError is the error returned in case of non-fatal parsing errors.
type TransceiverError struct {
	str string
}

func (e *TransceiverError) Error() string {
	return e.str
}

func newTransceiverError(format string, args ...interface{}) *TransceiverError {
	return &TransceiverError{
		str: fmt.Sprintf(format, args...),
	}
}

// TransceiverConf configures a Transceiver.
type TransceiverConf struct {
	// the reader from which frames will be read.
	Reader io.Reader
	// the writer to which frames will be written.
	Writer io.Writer

	// (optional) the dialect which contains the messages that will be encoded and decoded.
	// If not provided, messages are decoded in the MessageRaw struct.
	DialectDE *dialect.DecEncoder

	// (optional) the secret key used to validate incoming frames.
	// Non-signed frames are discarded. This feature requires v2 frames.
	InKey *frame.V2Key

	// Mavlink version used to encode messages. See Version
	// for the available options.
	OutVersion Version
	// the system id, added to every outgoing frame and used to identify this
	// node in the network.
	OutSystemId byte
	// (optional) the component id, added to every outgoing frame, defaults to 1.
	OutComponentId byte
	// (optional) the value to insert into the signature link id.
	// This feature requires v2 frames.
	OutSignatureLinkId byte
	// (optional) the secret key used to sign outgoing frames.
	// This feature requires v2 frames.
	OutKey *frame.V2Key
}

// Transceiver is a low-level Mavlink encoder and decoder that works with a Reader and a Writer.
type Transceiver struct {
	conf                 TransceiverConf
	readBuffer           *bufio.Reader
	writeBuffer          []byte
	curWriteSequenceId   byte
	curReadSignatureTime uint64
}

// New allocates a Transceiver, a low level frame encoder and decoder.
// See TransceiverConf for the options.
func New(conf TransceiverConf) (*Transceiver, error) {
	if conf.Reader == nil {
		return nil, fmt.Errorf("Reader not provided")
	}
	if conf.Writer == nil {
		return nil, fmt.Errorf("Writer not provided")
	}

	if conf.OutVersion == 0 {
		return nil, fmt.Errorf("OutVersion not provided")
	}
	if conf.OutSystemId < 1 {
		return nil, fmt.Errorf("SystemId must be >= 1")
	}
	if conf.OutComponentId < 1 {
		conf.OutComponentId = 1
	}
	if conf.OutKey != nil && conf.OutVersion != V2 {
		return nil, fmt.Errorf("OutKey requires V2 frames")
	}

	return &Transceiver{
		conf:        conf,
		readBuffer:  bufio.NewReaderSize(conf.Reader, bufferSize),
		writeBuffer: make([]byte, 0, bufferSize),
	}, nil
}

// Read reads a Frame from the reader.
// It must not be called by multiple routines in parallel.
func (p *Transceiver) Read() (frame.Frame, error) {
	magicByte, err := p.readBuffer.ReadByte()
	if err != nil {
		return nil, err
	}

	f, err := func() (frame.Frame, error) {
		switch magicByte {
		case frame.V1MagicByte:
			return &frame.V1Frame{}, nil

		case frame.V2MagicByte:
			return &frame.V2Frame{}, nil
		}

		return nil, newTransceiverError("invalid magic byte: %x", magicByte)
	}()
	if err != nil {
		return nil, err
	}

	err = f.Decode(p.readBuffer)
	if err != nil {
		return nil, newTransceiverError(err.Error())
	}

	if p.conf.InKey != nil {
		ff, ok := f.(*frame.V2Frame)
		if !ok {
			return nil, newTransceiverError("signature required but packet is not v2")
		}

		if sig := ff.GenSignature(p.conf.InKey); *sig != *ff.Signature {
			return nil, newTransceiverError("wrong signature")
		}

		// in UDP, packet order is not guaranteed. Therefore, we accept frames
		// with a timestamp within 10 seconds with respect to the previous frame.
		if p.curReadSignatureTime > 0 &&
			ff.SignatureTimestamp < (p.curReadSignatureTime-(10*100000)) {
			return nil, newTransceiverError("signature timestamp is too old")
		}

		if ff.SignatureTimestamp > p.curReadSignatureTime {
			p.curReadSignatureTime = ff.SignatureTimestamp
		}
	}

	// decode message if in dialect and validate checksum
	if p.conf.DialectDE != nil {
		if mp, ok := p.conf.DialectDE.MessageDEs[f.GetMessage().GetId()]; ok {
			if sum := f.GenChecksum(p.conf.DialectDE.MessageDEs[f.GetMessage().GetId()].CRCExtra()); sum != f.GetChecksum() {
				return nil, newTransceiverError("wrong checksum (expected %.4x, got %.4x, id=%d)",
					sum, f.GetChecksum(), f.GetMessage().GetId())
			}

			_, isV2 := f.(*frame.V2Frame)
			msg, err := mp.Decode(f.GetMessage().(*msg.MessageRaw).Content, isV2)
			if err != nil {
				return nil, newTransceiverError(err.Error())
			}

			switch ff := f.(type) {
			case *frame.V1Frame:
				ff.Message = msg
			case *frame.V2Frame:
				ff.Message = msg
			}
		}
	}

	return f, nil
}

// WriteMessage writes a Message into the writer.
// It must not be called by multiple routines in parallel.
func (p *Transceiver) WriteMessage(m msg.Message) error {
	var fr frame.Frame
	if p.conf.OutVersion == V1 {
		fr = &frame.V1Frame{Message: m}
	} else {
		fr = &frame.V2Frame{Message: m}
	}
	return p.writeFrameAndFill(fr)
}

func (p *Transceiver) writeFrameAndFill(fr frame.Frame) error {
	if fr.GetMessage() == nil {
		return fmt.Errorf("message is nil")
	}

	// do not touch the original frame, but work with a separate object
	// in such way that the frame can be encoded by other parsers in parallel
	safeFrame := fr.Clone()

	// fill SequenceId, SystemId, ComponentId
	switch ff := safeFrame.(type) {
	case *frame.V1Frame:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	case *frame.V2Frame:
		ff.SequenceId = p.curWriteSequenceId
		ff.SystemId = p.conf.OutSystemId
		ff.ComponentId = p.conf.OutComponentId
	}
	p.curWriteSequenceId++

	// fill CompatibilityFlag, IncompatibilityFlag if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok {
		ff.CompatibilityFlag = 0
		ff.IncompatibilityFlag = 0

		if p.conf.OutKey != nil {
			ff.IncompatibilityFlag |= frame.V2FlagSigned
		}
	}

	// encode message if it is not already encoded
	if _, ok := safeFrame.GetMessage().(*msg.MessageRaw); !ok {
		if p.conf.DialectDE == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.conf.DialectDE.MessageDEs[safeFrame.GetMessage().GetId()]
		if !ok {
			return fmt.Errorf("message cannot be encoded since it is not in the dialect")
		}

		_, isV2 := safeFrame.(*frame.V2Frame)
		byt, err := mp.Encode(safeFrame.GetMessage(), isV2)
		if err != nil {
			return err
		}

		msgRaw := &msg.MessageRaw{safeFrame.GetMessage().GetId(), byt}
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Message = msgRaw
		case *frame.V2Frame:
			ff.Message = msgRaw
		}

		// fill checksum
		switch ff := safeFrame.(type) {
		case *frame.V1Frame:
			ff.Checksum = ff.GenChecksum(p.conf.DialectDE.MessageDEs[ff.GetMessage().GetId()].CRCExtra())
		case *frame.V2Frame:
			ff.Checksum = ff.GenChecksum(p.conf.DialectDE.MessageDEs[ff.GetMessage().GetId()].CRCExtra())
		}
	}

	// fill SignatureLinkId, SignatureTimestamp, Signature if v2
	if ff, ok := safeFrame.(*frame.V2Frame); ok && p.conf.OutKey != nil {
		ff.SignatureLinkId = p.conf.OutSignatureLinkId
		// Timestamp in 10 microsecond units since 1st January 2015 GMT time
		ff.SignatureTimestamp = uint64(time.Since(signatureReferenceDate)) / 10000
		ff.Signature = ff.GenSignature(p.conf.OutKey)
	}

	return p.WriteFrame(safeFrame)
}

// WriteFrame writes a Frame into the writer.
// It must not be called by multiple routines in parallel.
// This function is intended only for routing pre-existing frames to other nodes,
// since all frame fields must be filled manually.
func (p *Transceiver) WriteFrame(fr frame.Frame) error {
	m := fr.GetMessage()
	if m == nil {
		return fmt.Errorf("message is nil")
	}

	// encode message if it is not already encoded
	if _, ok := m.(*msg.MessageRaw); !ok {
		if p.conf.DialectDE == nil {
			return fmt.Errorf("message cannot be encoded since dialect is nil")
		}

		mp, ok := p.conf.DialectDE.MessageDEs[m.GetId()]
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
		m = &msg.MessageRaw{m.GetId(), byt}
	}

	buf, err := fr.Encode(p.writeBuffer, m.(*msg.MessageRaw).Content)
	if err != nil {
		return err
	}

	// do not check n, since io.Writer is not allowed to return n < len(buf)
	// without throwing an error
	_, err = p.conf.Writer.Write(buf)
	return err
}
