package gomavlib

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParserNil(t *testing.T) {
	parser, err := NewParser(ParserConf{
		Reader:         bytes.NewReader(nil),
		Writer:         bytes.NewBuffer(nil),
		Dialect:        nil,
		OutSystemId:    1,
	})
	require.NoError(t, err)
	frame := &FrameV1{ Message: nil }
	err = parser.Write(frame, true)
	require.Error(t, err)
}

type MessageTest5 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest5) GetId() uint32 {
	return 5
}

type MessageTest6 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest6) GetId() uint32 {
	return 0x0607
}

type MessageTest8 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest8) GetId() uint32 {
	return 8
}

var testDialect = MustDialect(3, []Message{
	&MessageTest5{},
	&MessageTest6{},
	&MessageTest8{},
})

func testFrameDecode(t *testing.T, dialect *Dialect, key *FrameSignatureKey, byts [][]byte, frames []Frame) {
	for i, byt := range byts {
		parser, err := NewParser(ParserConf{
			Reader:         bytes.NewReader(byt),
			Writer:         bytes.NewBuffer(nil),
			Dialect:        dialect,
			OutSystemId:    1,
			InSignatureKey: key,
		})
		require.NoError(t, err)
		frame, err := parser.Read()
		require.NoError(t, err)
		require.Equal(t, frames[i], frame)
	}
}

func testFrameEncode(t *testing.T, dialect *Dialect, key *FrameSignatureKey, byts [][]byte, frames []Frame) {
	for i, frame := range frames {
		buf := bytes.NewBuffer(nil)
		parser, err := NewParser(ParserConf{
			Reader:      bytes.NewBuffer(nil),
			Writer:      buf,
			OutSystemId: 1,
			Dialect:     dialect,
		})
		require.NoError(t, err)
		err = parser.Write(frame, true)
		require.NoError(t, err)
		require.Equal(t, byts[i], buf.Bytes())
	}
}

var testFpV1Bytes = [][]byte{
	[]byte("\xFE\x05\x27\x01\x02\x05\x10\x10\x10\x10\x10\xe5\x66"),
	[]byte("\xFE\x05\x27\x01\x02\x08\x10\x10\x10\x10\x10\xfa\xc7"),
}

var testFpV1Frames = []Frame{
	&FrameV1{
		SequenceId:  0x27,
		SystemId:    0x01,
		ComponentId: 0x02,
		Message: &MessageRaw{
			0x05,
			[]byte("\x10\x10\x10\x10\x10"),
		},
		Checksum: 0x66e5,
	},
	&FrameV1{
		SequenceId:  0x27,
		SystemId:    0x01,
		ComponentId: 0x02,
		Message: &MessageRaw{
			0x08,
			[]byte("\x10\x10\x10\x10\x10"),
		},
		Checksum: 0xc7fa,
	},
}

var testFpV1FramesDialect = []Frame{
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
	&FrameV1{
		SequenceId:  0x27,
		SystemId:    0x01,
		ComponentId: 0x02,
		Message: &MessageTest8{
			'\x10',
			binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
		},
		Checksum: 0xc7fa,
	},
}

func TestParserV1RawDec(t *testing.T) {
	testFrameDecode(t, nil, nil, testFpV1Bytes, testFpV1Frames)
}

func TestParserV1RawEnc(t *testing.T) {
	testFrameEncode(t, nil, nil, testFpV1Bytes, testFpV1Frames)
}

func TestParserV1DialectDec(t *testing.T) {
	testFrameDecode(t, testDialect, nil, testFpV1Bytes, testFpV1FramesDialect)
}

func TestParserV1DialectEnc(t *testing.T) {
	testFrameEncode(t, testDialect, nil, testFpV1Bytes, testFpV1FramesDialect)
}

var testFpV2Bytes = [][]byte{
	[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
}

var testFpV2Frames = []Frame{
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
}

var testFpV2FramesDialect = []Frame{
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
}

func TestParserV2RawDec(t *testing.T) {
	testFrameDecode(t, nil, nil, testFpV2Bytes, testFpV2Frames)
}

func TestParserV2RawEnc(t *testing.T) {
	testFrameEncode(t, nil, nil, testFpV2Bytes, testFpV2Frames)
}

func TestParserV2DialectDec(t *testing.T) {
	testFrameDecode(t, testDialect, nil, testFpV2Bytes, testFpV2FramesDialect)
}

func TestParserV2DialectEnc(t *testing.T) {
	testFrameEncode(t, testDialect, nil, testFpV2Bytes, testFpV2FramesDialect)
}

/* Test vectors generated with

( docker build - -t temp << EOF
FROM amd64/python:3-stretch
RUN apt update && apt install -y --no-install-recommends \
    git \
    gcc \
    python3-dev \
    python3-setuptools \
    python3-wheel \
    python3-pip \
    python3-future \
    python3-lxml \
    && pip3 install pymavlink
EOF
) && docker run --rm -it temp python3 -c \
"import io; from pymavlink.dialects.v20 import ardupilotmega; \
f = io.BytesIO(); \
mav = ardupilotmega.MAVLink(f); \
mav.signing.secret_key = (chr(79)*32).encode(); \
mav.signing.link_id = 1; \
mav.signing.timestamp = 2; \
mav.signing.sign_outgoing = True; \
mav.heartbeat_send(type=1, autopilot=2, base_mode=3, custom_mode=4, system_status=5); \
f.seek(0); \
print(''.join('\\\x{:02x}'.format(x) for x in f.read()));"

( docker build - -t temp << EOF
FROM amd64/python:3-stretch
RUN apt update && apt install -y --no-install-recommends \
    git \
    gcc \
    python3-dev \
    python3-setuptools \
    python3-wheel \
    python3-pip \
    python3-future \
    python3-lxml \
    && pip3 install pymavlink
EOF
) && docker run --rm -it temp python3 -c \
"import io; from pymavlink.dialects.v20 import ardupilotmega; \
f = io.BytesIO(); \
mav = ardupilotmega.MAVLink(f); \
mav.signing.secret_key = (chr(79)*32).encode(); \
mav.signing.link_id = 3; \
mav.signing.timestamp = 4; \
mav.signing.sign_outgoing = True; \
mav.optical_flow_send(time_usec=1, sensor_id=2, flow_x=3, flow_y=4, flow_comp_m_x=5, flow_comp_m_y=6, quality=7, ground_distance=8, flow_rate_y=1); \
f.seek(0); \
print(''.join('\\\x{:02x}'.format(x) for x in f.read()));"

*/

var testFpV2SigKey = NewFrameSignatureKey(bytes.Repeat([]byte("\x4F"), 32))

var testFpV2SigDialect = MustDialect(3, []Message{
	&MessageHeartbeat{},
	&MessageOpticalFlow{},
})

var testFpV2SigBytes = [][]byte{
	[]byte("\xfd\x09\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0e\x47\x04\x0c\xef\x9b"),
	[]byte("\xfd\x22\x01\x00\x00\x00\x00\x64\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa0\x40\x00\x00\xc0\x40\x00\x00\x00\x41\x03\x00\x04\x00\x02\x07\x00\x00\x00\x00\x00\x00\x80\x3f\x77\xfb\x03\x04\x00\x00\x00\x00\x00\xa8\x88\x09\x39\xb2\x60"),
}

var testFpV2SigFrames = []Frame{
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
		Signature:          &FrameSignature{0x0e, 0x47, 0x04, 0x0c, 0xef, 0x9b},
	},
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
		Signature:          &FrameSignature{0xa8, 0x88, 0x9, 0x39, 0xb2, 0x60},
	},
}

func TestParserV2SignatureDec(t *testing.T) {
	testFrameDecode(t, testFpV2SigDialect, testFpV2SigKey, testFpV2SigBytes, testFpV2SigFrames)
}

func TestParserV2SignatureEnc(t *testing.T) {
	testFrameEncode(t, testFpV2SigDialect, testFpV2SigKey, testFpV2SigBytes, testFpV2SigFrames)
}
