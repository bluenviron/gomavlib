package gomavlib

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

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

var testFpV1Frames = []*FrameV1{
	{
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

var testFpV1FramesDialect = []*FrameV1{
	{
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

func TestFrameParserV1RawDec(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{})
	for i, byt := range testFpV1Bytes {
		frame, err := parser.Decode(byt, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(frame, testFpV1Frames[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", frame, testFpV1Frames[i])
		}
	}
}

func TestFrameParserV1RawEnc(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{})
	for i, frame := range testFpV1Frames {
		byt, err := parser.Encode(frame, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(byt, testFpV1Bytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testFpV1Bytes[i])
		}
	}
}

func TestFrameParserV1DialectDec(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{Dialect: testDialect})
	for i, byt := range testFpV1Bytes {
		frame, err := parser.Decode(byt, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(frame, testFpV1FramesDialect[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", frame, testFpV1FramesDialect[i])
		}
	}
}

func TestFrameParserV1DialectEnc(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{Dialect: testDialect})
	for i, frame := range testFpV1FramesDialect {
		byt, err := parser.Encode(frame, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(byt, testFpV1Bytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testFpV1Bytes[i])
		}
	}
}

var testFpV2Bytes = [][]byte{
	[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
}

var testFpV2Frames = []*FrameV2{
	{
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

var testFpV2FramesDialect = []*FrameV2{
	{
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

func TestFrameParserV2RawDec(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{})
	for i, byt := range testFpV2Bytes {
		frame, err := parser.Decode(byt, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(frame, testFpV2Frames[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", frame, testFpV2Frames[i])
		}
	}
}

func TestFrameParserV2RawEnc(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{})
	for i, frame := range testFpV2Frames {
		byt, err := parser.Encode(frame, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(byt, testFpV2Bytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testFpV2Bytes[i])
		}
	}
}

func TestFrameParserV2DialectDec(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{Dialect: testDialect})
	for i, byt := range testFpV2Bytes {
		frame, err := parser.Decode(byt, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(frame, testFpV2FramesDialect[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", frame, testFpV2FramesDialect[i])
		}
	}
}

func TestFrameParserV2DialectEnc(t *testing.T) {
	parser, _ := NewFrameParser(FrameParserConf{Dialect: testDialect})
	for i, frame := range testFpV2FramesDialect {
		byt, err := parser.Encode(frame, true, nil)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(byt, testFpV2Bytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testFpV2Bytes[i])
		}
	}
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

var testFpV2SigFrames = []*FrameV2{
	{
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
		Signature:          Signature{14, 71, 4, 12, 239, 155},
	},
}

func TestFrameParserV2SignatureDec(t *testing.T) {
	var key SignatureKey
	copy(key[:], bytes.Repeat([]byte("\x4F"), 32))

	parser, _ := NewFrameParser(FrameParserConf{
		Dialect: []Message{&MessageHeartbeat{}},
	})

	for i, byt := range testFpV2SigBytes {
		frame, err := parser.Decode(byt, true, &key)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(frame, testFpV2SigFrames[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", frame, testFpV2SigFrames[i])
		}
	}
}

func TestFrameParserV2SignatureEnc(t *testing.T) {
	var key SignatureKey
	copy(key[:], bytes.Repeat([]byte("\x4F"), 32))

	parser, _ := NewFrameParser(FrameParserConf{
		Dialect: []Message{&MessageHeartbeat{}},
	})

	for i, frame := range testFpV2SigFrames {
		byt, err := parser.Encode(frame, true, &key)
		if err != nil {
			t.Fatal(err)
		}
		if reflect.DeepEqual(byt, testFpV2SigBytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testFpV2SigBytes[i])
		}
	}
}
