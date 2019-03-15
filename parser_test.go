package gomavlib

import (
	"bytes"
	"encoding/binary"
	"github.com/stretchr/testify/require"
	"testing"
)

func testFrameDecode(t *testing.T, dialect []Message, key *FrameSignatureKey, byts [][]byte, frames []Frame) {
	parser, _ := NewParser(ParserConf{Dialect: dialect})
	for i, byt := range byts {
		frame, err := parser.Decode(byt, true, key)
		require.NoError(t, err)
		require.Equal(t, frame, frames[i])
	}
}

func testFrameEncode(t *testing.T, dialect []Message, key *FrameSignatureKey, byts [][]byte, frames []Frame) {
	parser, _ := NewParser(ParserConf{Dialect: dialect})
	for i, frame := range frames {
		byt, err := parser.Encode(frame, true, key)
		require.NoError(t, err)
		require.Equal(t, byt, byts[i])
	}
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

var testDialect = []Message{
	&MessageTest5{},
	&MessageTest6{},
	&MessageTest8{},
}

var testFpV1Bytes = [][]byte{
	[]byte("\xFE\x05\x27\x01\x02\x05\x10\x10\x10\x10\x10\xe5\x66"),
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
print(f.read());"

*/

var testFpV2SigBytes = [][]byte{
	[]byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b"),
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
		Checksum:           53721,
		SignatureLinkId:    1,
		SignatureTimestamp: 2,
		Signature:          &FrameSignature{14, 71, 4, 12, 239, 155},
	},
}

var testFpV2Key = NewFrameSignatureKey(bytes.Repeat([]byte("\x4F"), 32))

func TestParserV2FrameSignatureDec(t *testing.T) {
	testFrameDecode(t, []Message{&MessageHeartbeat{}}, testFpV2Key, testFpV2SigBytes, testFpV2SigFrames)
}

func TestParserV2FrameSignatureEnc(t *testing.T) {
	testFrameEncode(t, []Message{&MessageHeartbeat{}}, testFpV2Key, testFpV2SigBytes, testFpV2SigFrames)
}
