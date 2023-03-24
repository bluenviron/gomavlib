package frame

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v2/pkg/dialect"
	"github.com/bluenviron/gomavlib/v2/pkg/message"
)

type (
	MAV_TYPE      uint32 //nolint:revive
	MAV_AUTOPILOT uint32 //nolint:revive
	MAV_MODE_FLAG uint32 //nolint:revive
	MAV_STATE     uint32 //nolint:revive
)

type MessageTest5 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest5) GetID() uint32 {
	return 5
}

type MessageTest6 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest6) GetID() uint32 {
	return 0x0607
}

type MessageTest8 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest8) GetID() uint32 {
	return 8
}

type MessageTest9 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest9) GetID() uint32 {
	return 300
}

type MessageHeartbeat struct {
	Type           MAV_TYPE      `mavenum:"uint8"`
	Autopilot      MAV_AUTOPILOT `mavenum:"uint8"`
	BaseMode       MAV_MODE_FLAG `mavenum:"uint8"`
	CustomMode     uint32
	SystemStatus   MAV_STATE `mavenum:"uint8"`
	MavlinkVersion uint8
}

func (*MessageHeartbeat) GetID() uint32 {
	return 0
}

type MessageOpticalFlow struct {
	TimeUsec       uint64
	SensorId       uint8 //nolint:revive
	FlowX          int16
	FlowY          int16
	FlowCompMX     float32
	FlowCompMY     float32
	Quality        uint8
	GroundDistance float32
	FlowRateX      float32 `mavext:"true"`
	FlowRateY      float32 `mavext:"true"`
}

func (*MessageOpticalFlow) GetID() uint32 {
	return 100
}

var testDialectRW = func() *dialect.ReadWriter {
	d := &dialect.Dialect{3, []message.Message{ //nolint:govet
		&MessageTest5{},
		&MessageTest6{},
		&MessageTest9{},
		&MessageHeartbeat{},
		&MessageOpticalFlow{},
	}}
	de, err := dialect.NewReadWriter(d)
	if err != nil {
		panic(err)
	}
	return de
}()

var casesReadWrite = []struct {
	name      string
	dialectRW *dialect.ReadWriter
	key       *V2Key
	frame     Frame
	raw       []byte
}{
	{
		"v1 frame with nil content",
		nil,
		nil,
		&V1Frame{
			SequenceID:  0x01,
			SystemID:    0x02,
			ComponentID: 0x03,
			Message: &message.MessageRaw{
				ID:      4,
				Payload: nil,
			},
			Checksum: 0x0807,
		},
		[]byte("\xFE\x00\x01\x02\x03\x04\x07\x08"),
	},
	{
		"v1 frame with encoded message",
		nil,
		nil,
		&V1Frame{
			SequenceID:  0x27,
			SystemID:    0x01,
			ComponentID: 0x02,
			Message: &message.MessageRaw{ //nolint:govet
				8,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0xc7fa,
		},
		[]byte("\xFE\x05\x27\x01\x02\x08\x10\x10\x10\x10\x10\xfa\xc7"),
	},
	{
		"v1 frame with decoded message",
		testDialectRW,
		nil,
		&V1Frame{
			SequenceID:  0x27,
			SystemID:    0x01,
			ComponentID: 0x02,
			Message: &MessageTest5{
				'\x10',
				binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
			},
			Checksum: 0x66e5,
		},
		[]byte("\xFE\x05\x27\x01\x02\x05\x10\x10\x10\x10\x10\xe5\x66"),
	},
	{
		"v2 frame with nil content",
		testDialectRW,
		nil,
		&V2Frame{
			IncompatibilityFlag: 0,
			CompatibilityFlag:   0,
			SequenceID:          3,
			SystemID:            4,
			ComponentID:         5,
			Message: &message.MessageRaw{ //nolint:govet
				4,
				nil,
			},
			Checksum: 0x0ab7,
		},
		[]byte("\xFD\x00\x00\x00\x03\x04\x05\x04\x00\x00\xb7\x0a"),
	},
	{
		"v2 frame with encoded message",
		nil,
		nil,
		&V2Frame{
			IncompatibilityFlag: 0x00,
			CompatibilityFlag:   0x00,
			SequenceID:          0x8F,
			SystemID:            0x01,
			ComponentID:         0x02,
			Message: &message.MessageRaw{ //nolint:govet
				0x0607,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0x0349,
		},
		[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
	},
	{
		"v2 frame with decoded message",
		testDialectRW,
		nil,
		&V2Frame{
			IncompatibilityFlag: 0,
			CompatibilityFlag:   0x00,
			SequenceID:          0x8F,
			SystemID:            0x01,
			ComponentID:         0x02,
			Message: &MessageTest6{
				'\x10',
				binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
			},
			Checksum: 0x0349,
		},
		[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
	},
	{
		"v2 frame with decoded message, signed",
		testDialectRW,
		NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
		&V2Frame{
			IncompatibilityFlag: 0x01,
			CompatibilityFlag:   0x00,
			SequenceID:          0x00,
			SystemID:            0x00,
			ComponentID:         0x00,
			Message: &MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     4,
				SystemStatus:   5,
				MavlinkVersion: 3,
			},
			Checksum:           0xd1d9,
			SignatureLinkID:    1,
			SignatureTimestamp: 2,
			Signature:          &V2Signature{0x0e, 0x47, 0x04, 0x0c, 0xef, 0x9b},
		},
		[]byte("\xFD\x09\x01\x00\x00\x00\x00\x00" +
			"\x00\x00\x04\x00\x00\x00\x01\x02" +
			"\x03\x05\x03\xd9\xd1\x01\x02\x00" +
			"\x00\x00\x00\x00\x0e\x47\x04\x0c\xef\x9b"),
	},
	{
		"v2 frame with decoded message, signed",
		testDialectRW,
		NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
		&V2Frame{
			IncompatibilityFlag: 0x01,
			CompatibilityFlag:   0x00,
			SequenceID:          0x00,
			SystemID:            0x00,
			ComponentID:         0x00,
			Message: &MessageOpticalFlow{
				TimeUsec:       1,
				SensorId:       2,
				FlowX:          3,
				FlowY:          4,
				FlowCompMX:     5,
				FlowCompMY:     6,
				Quality:        7,
				GroundDistance: 8,
				FlowRateY:      1,
			},
			Checksum:           0xfb77,
			SignatureLinkID:    3,
			SignatureTimestamp: 4,
			Signature:          &V2Signature{0xa8, 0x88, 0x9, 0x39, 0xb2, 0x60},
		},
		[]byte("\xFD\x22\x01\x00\x00\x00\x00\x64" +
			"\x00\x00\x01\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x00\xa0\x40\x00\x00\xc0\x40" +
			"\x00\x00\x00\x41\x03\x00\x04\x00\x02" +
			"\x07\x00\x00\x00\x00\x00\x00\x80\x3f" +
			"\x77\xfb\x03\x04\x00\x00\x00\x00\x00" +
			"\xa8\x88\x09\x39\xb2\x60"),
	},
}

func TestReaderNewErrors(t *testing.T) {
	_, err := NewReader(ReaderConf{})
	require.EqualError(t, err, "Reader not provided")
}

func TestReader(t *testing.T) {
	for _, ca := range casesReadWrite {
		t.Run(ca.name, func(t *testing.T) {
			reader, err := NewReader(ReaderConf{
				Reader:    bytes.NewReader(ca.raw),
				DialectRW: ca.dialectRW,
				InKey:     ca.key,
			})
			require.NoError(t, err)

			frame, err := reader.Read()
			require.NoError(t, err)
			require.Equal(t, ca.frame, frame)
		})
	}
}

func TestReaderErrors(t *testing.T) {
	for _, ca := range []struct {
		name      string
		dialectRW *dialect.ReadWriter
		key       *V2Key
		byts      []byte
		err       string
	}{
		{
			"empty",
			nil,
			nil,
			nil,
			"EOF",
		},
		{
			"invalid magic byte",
			nil,
			nil,
			[]byte{0x07},
			"invalid magic byte: 7",
		},
		{
			"invalid frame",
			nil,
			nil,
			[]byte{0xFE},
			"EOF",
		},
		{
			"v1 frame but signature required",
			nil,
			NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
			[]byte{
				0xFE, 0x05, 0x27, 0x01, 0x02, 0x08, 0x10, 0x10,
				0x10, 0x10, 0x10, 0xfa, 0xc7,
			},
			"signature required but packet is not v2",
		},
		{
			"v2 frame unknown incompatibility flag",
			nil,
			nil,
			[]byte{
				0xfd, 0x5, 0x4, 0x0, 0x8f, 0x1, 0x2, 0x7,
				0x6, 0x0, 0x10, 0x10, 0x10, 0x10, 0x10, 0x49,
				0x3,
			},
			"unknown incompatibility flag: 4",
		},
		{
			"wrong signature",
			nil,
			NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
			[]byte{
				0xFD, 0x09, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x02,
				0x03, 0x05, 0x03, 0xd9, 0xd1, 0x01, 0x02, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x0e, 0x47, 0x04, 0x0c,
				0xef, 0x9c,
			},
			"wrong signature",
		},
		{
			"wrong checksum",
			testDialectRW,
			nil,
			[]byte{
				0xFE, 0x05, 0x27, 0x01, 0x02, 0x05, 0x10, 0x10,
				0x10, 0x10, 0x10, 0xe6, 0x66,
			},
			"wrong checksum, expected 66e5, got 66e6, message id is 5",
		},
		{
			"message read error",
			nil,
			nil,
			[]byte{254, 3, 39, 1, 2, 5, 168, 233},
			"unexpected EOF",
		},
		{
			"message decode error",
			testDialectRW,
			nil,
			[]byte{254, 0, 39, 1, 2, 5, 168, 233},
			"unable to decode message: wrong size: expected 5, got 0",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			reader, err := NewReader(ReaderConf{
				Reader:    bytes.NewReader(ca.byts),
				DialectRW: ca.dialectRW,
				InKey:     ca.key,
			})
			require.NoError(t, err)

			_, err = reader.Read()
			require.EqualError(t, err, ca.err)
		})
	}
}

func TestReaderErrorSignatureTimestamp(t *testing.T) {
	var buf bytes.Buffer

	msgByts := []byte{4, 0, 0, 0, 1, 2, 3, 5, 3}

	f := &V2Frame{
		IncompatibilityFlag: 0x01,
		CompatibilityFlag:   0x00,
		SequenceID:          0x00,
		SystemID:            0x00,
		ComponentID:         0x00,
		Message: &message.MessageRaw{
			ID:      0,
			Payload: msgByts,
		},
		Checksum:           0xd1d9,
		SignatureLinkID:    1,
		SignatureTimestamp: 20000000,
	}
	f.Signature = f.GenerateSignature(NewV2Key(bytes.Repeat([]byte("\x4F"), 32)))
	buf2 := make([]byte, 1024)
	n, err := f.encodeTo(buf2, msgByts)
	require.NoError(t, err)
	buf.Write(buf2[:n])

	f = &V2Frame{
		IncompatibilityFlag: 0x01,
		CompatibilityFlag:   0x00,
		SequenceID:          0x00,
		SystemID:            0x00,
		ComponentID:         0x00,
		Message: &message.MessageRaw{
			ID:      0,
			Payload: msgByts,
		},
		Checksum:           0xd1d9,
		SignatureLinkID:    1,
		SignatureTimestamp: 2,
	}
	f.Signature = f.GenerateSignature(NewV2Key(bytes.Repeat([]byte("\x4F"), 32)))
	buf2 = make([]byte, 1024)
	n, err = f.encodeTo(buf2, msgByts)
	require.NoError(t, err)
	buf.Write(buf2[:n])

	reader, err := NewReader(ReaderConf{
		Reader: &buf,
		InKey:  NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
	})
	require.NoError(t, err)

	_, err = reader.Read()
	require.NoError(t, err)

	_, err = reader.Read()
	require.EqualError(t, err, "signature timestamp is too old")
}
