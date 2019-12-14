package gomavlib

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"
)

var casesParser = []struct {
	name    string
	dialect *Dialect
	key     *Key
	frame   Frame
	raw     []byte
}{
	{
		"v1 frame with nil content",
		nil,
		nil,
		&FrameV1{
			SequenceId:  0x01,
			SystemId:    0x02,
			ComponentId: 0x03,
			Message: &MessageRaw{
				Id:      4,
				Content: nil,
			},
			Checksum: 0x0807,
		},
		[]byte("\xFE\x00\x01\x02\x03\x04\x07\x08"),
	},
	{
		"v1 frame with encoded message",
		nil,
		nil,
		&FrameV1{
			SequenceId:  0x27,
			SystemId:    0x01,
			ComponentId: 0x02,
			Message: &MessageRaw{
				8,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0xc7fa,
		},
		[]byte("\xFE\x05\x27\x01\x02\x08\x10\x10\x10\x10\x10\xfa\xc7"),
	},
	{
		"v1 frame with decoded message",
		testDialect,
		nil,
		&FrameV1{
			SequenceId:  0x27,
			SystemId:    0x01,
			ComponentId: 0x02,
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
		testDialect,
		nil,
		&FrameV2{
			IncompatibilityFlag: 0,
			CompatibilityFlag:   0,
			SequenceId:          3,
			SystemId:            4,
			ComponentId:         5,
			Message: &MessageRaw{
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
		&FrameV2{
			IncompatibilityFlag: 0x00,
			CompatibilityFlag:   0x00,
			SequenceId:          0x8F,
			SystemId:            0x01,
			ComponentId:         0x02,
			Message: &MessageRaw{
				0x0607,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0x0349,
		},
		[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
	},
	{
		"v2 frame with decoded message",
		testDialect,
		nil,
		&FrameV2{
			IncompatibilityFlag: 0x00,
			CompatibilityFlag:   0x00,
			SequenceId:          0x8F,
			SystemId:            0x01,
			ComponentId:         0x02,
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
		testDialect,
		NewKey(bytes.Repeat([]byte("\x4F"), 32)),
		&FrameV2{
			IncompatibilityFlag: 0x01,
			CompatibilityFlag:   0x00,
			SequenceId:          0x00,
			SystemId:            0x00,
			ComponentId:         0x00,
			Message: &MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     4,
				SystemStatus:   5,
				MavlinkVersion: 3,
			},
			Checksum:           0xd1d9,
			SignatureLinkId:    1,
			SignatureTimestamp: 2,
			Signature:          &Signature{0x0e, 0x47, 0x04, 0x0c, 0xef, 0x9b},
		},
		[]byte("\xFD\x09\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0e\x47\x04\x0c\xef\x9b"),
	},
	{
		"v2 frame with decoded message, signed",
		testDialect,
		NewKey(bytes.Repeat([]byte("\x4F"), 32)),
		&FrameV2{
			IncompatibilityFlag: 0x01,
			CompatibilityFlag:   0x00,
			SequenceId:          0x00,
			SystemId:            0x00,
			ComponentId:         0x00,
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
			SignatureLinkId:    3,
			SignatureTimestamp: 4,
			Signature:          &Signature{0xa8, 0x88, 0x9, 0x39, 0xb2, 0x60},
		},
		[]byte("\xFD\x22\x01\x00\x00\x00\x00\x64\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa0\x40\x00\x00\xc0\x40\x00\x00\x00\x41\x03\x00\x04\x00\x02\x07\x00\x00\x00\x00\x00\x00\x80\x3f\x77\xfb\x03\x04\x00\x00\x00\x00\x00\xa8\x88\x09\x39\xb2\x60"),
	},
}

func TestParserDecode(t *testing.T) {
	for _, c := range casesParser {
		t.Run(c.name, func(t *testing.T) {
			parser, err := NewParser(ParserConf{
				Reader:      bytes.NewReader(c.raw),
				Writer:      bytes.NewBuffer(nil),
				Dialect:     c.dialect,
				OutSystemId: 1,
				InKey:       c.key,
			})
			require.NoError(t, err)
			frame, err := parser.Read()
			require.NoError(t, err)
			require.Equal(t, c.frame, frame)
		})
	}
}

func TestParserEncode(t *testing.T) {
	for _, c := range casesParser {
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			parser, err := NewParser(ParserConf{
				Reader:      bytes.NewBuffer(nil),
				Writer:      buf,
				OutSystemId: 1,
				Dialect:     c.dialect,
			})
			require.NoError(t, err)
			err = parser.WriteFrameRaw(c.frame)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
		})
	}
}

var casesParserWriteFrameAndFill = []struct {
	name  string
	frame Frame
	raw   []byte
}{
	{
		"v1 frame",
		&FrameV1{
			Message: &MessageTest5{
				'\x10',
				binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
			},
		},
		[]byte("\xFE\x05\x00\x01\x01\x05\x10\x10\x10\x10\x10\x75\x84"),
	},
	{
		"v1 frame, again",
		&FrameV1{
			Message: &MessageTest5{
				'\x10',
				binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
			},
		},
		[]byte("\xFE\x05\x01\x01\x01\x05\x10\x10\x10\x10\x10\x52\xA8"),
	},
	{
		"v2 frame with decoded message, signed",
		&FrameV2{
			Message: &MessageHeartbeat{
				Type:           1,
				Autopilot:      2,
				BaseMode:       3,
				CustomMode:     4,
				SystemStatus:   5,
				MavlinkVersion: 3,
			},
		},
		[]byte("\xFD\x09\x01\x00\x02\x01\x01\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\x28\xf3\x00\xe0\xf8\xd4\xb6\x8e\x0c\xff\x39\x30\x7f\xf6\x2e"),
	},
}

func TestParserWriteFrameAndFill(t *testing.T) {
	// fake current time in order to obtain deterministic signatures
	wayback := time.Date(2019, time.May, 18, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	// use a single parser for all tests
	buf := bytes.NewBuffer(nil)
	parser, err := NewParser(ParserConf{
		Reader:      bytes.NewBuffer(nil),
		Writer:      buf,
		OutSystemId: 1,
		Dialect:     testDialect,
		OutKey:      NewKey(bytes.Repeat([]byte("\x4F"), 32)),
	})
	require.NoError(t, err)

	for _, c := range casesParserWriteFrameAndFill {
		t.Run(c.name, func(t *testing.T) {
			err = parser.WriteFrameAndFill(c.frame)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
			buf.Next(len(c.raw))
		})
	}
}

func TestParserEncodeNilMsg(t *testing.T) {
	parser, err := NewParser(ParserConf{
		Reader:      bytes.NewReader(nil),
		Writer:      bytes.NewBuffer(nil),
		Dialect:     nil,
		OutSystemId: 1,
	})
	require.NoError(t, err)

	frame := &FrameV1{Message: nil}
	err = parser.WriteFrameRaw(frame)
	require.Error(t, err)
}

// ensure that the Frame is left untouched by the encoding procedure, and
// therefore WriteFrameAndFill() and WriteFrameRaw() can be called by multiple routines in parallel
func TestParserWriteFrameIsConst(t *testing.T) {
	parser, err := NewParser(ParserConf{
		Reader:      bytes.NewReader(nil),
		Writer:      bytes.NewBuffer(nil),
		Dialect:     MustDialect(3, []Message{&MessageHeartbeat{}}),
		OutSystemId: 1,
		OutKey:      NewKey(bytes.Repeat([]byte("\x7C"), 32)),
	})
	require.NoError(t, err)

	frame := &FrameV2{
		Message: &MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     4,
			SystemStatus:   5,
			MavlinkVersion: 3,
		},
	}
	original := frame.Clone()

	err = parser.WriteFrameAndFill(frame)
	require.NoError(t, err)
	require.Equal(t, frame, original)

	err = parser.WriteFrameRaw(frame)
	require.NoError(t, err)
	require.Equal(t, frame, original)
}
