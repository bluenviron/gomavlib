package gomavlib

import (
	"bytes"
	"reflect"
	"testing"
)

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
"from pymavlink.dialects.v20 import common; print(common.MAVLink_optical_flow_message.crc_extra);"

*/

type MessageHeartbeat struct {
	Type           uint8
	Autopilot      uint8
	BaseMode       uint8
	CustomMode     uint32
	SystemStatus   uint8
	MavlinkVersion uint8
}

func (m *MessageHeartbeat) GetId() uint32 {
	return 0
}

type MessageSysStatus struct {
	OnboardControlSensorsPresent uint32
	OnboardControlSensorsEnabled uint32
	OnboardControlSensorsHealth  uint32
	Load                         uint16
	VoltageBattery               uint16
	CurrentBattery               int16
	BatteryRemaining             int8
	DropRateComm                 uint16
	ErrorsComm                   uint16
	ErrorsCount1                 uint16
	ErrorsCount2                 uint16
	ErrorsCount3                 uint16
	ErrorsCount4                 uint16
}

func (m *MessageSysStatus) GetId() uint32 {
	return 1
}

type MessageChangeOperatorControl struct {
	TargetSystem   uint8
	ControlRequest uint8
	Version        uint8
	Passkey        string `mavlen:"25"`
}

func (m *MessageChangeOperatorControl) GetId() uint32 {
	return 5
}

type MessageAttitudeQuaternionCov struct {
	TimeUsec   uint64
	Q          [4]float32
	Rollspeed  float32
	Pitchspeed float32
	Yawspeed   float32
	Covariance [9]float32
}

func (m *MessageAttitudeQuaternionCov) GetId() uint32 {
	return 61
}

type MessageOpticalFlow struct {
	TimeUsec       uint64
	SensorId       uint8
	FlowX          int16
	FlowY          int16
	FlowCompMX     float32
	FlowCompMY     float32
	Quality        uint8
	GroundDistance float32
	FlowRateX      float32 `mavext:"true"`
	FlowRateY      float32 `mavext:"true"`
}

func (*MessageOpticalFlow) GetId() uint32 {
	return 100
}

func TestMessageParserCRC(t *testing.T) {
	var ins = []Message{
		&MessageHeartbeat{},
		&MessageSysStatus{},
		&MessageChangeOperatorControl{},
		&MessageAttitudeQuaternionCov{},
		&MessageOpticalFlow{},
	}
	var outs = []byte{
		50,
		124,
		217,
		167,
		175,
	}
	for i, in := range ins {
		mp, _ := newMessageParser(in)
		if mp.crcExtra != outs[i] {
			t.Fatalf("wrong crc extra: %v vs %v", mp.crcExtra, outs[i])
		}
	}
}

var testMpV1Bytes = [][]byte{
	[]byte("\x06\x00\x00\x00\x01\x02\x03\x04\x05"),
	bytes.Repeat([]byte("\x01"), 31),
	[]byte("\x01\x01\x01\x74\x65\x73\x74\x69\x6e\x67\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
	append([]byte("\x02\x00\x00\x00\x00\x00\x00\x00"), bytes.Repeat([]byte("\x00\x00\x80\x3F"), 16)...),
}

var testMpV1Parsers = []Message{
	&MessageHeartbeat{},
	&MessageSysStatus{},
	&MessageChangeOperatorControl{},
	&MessageAttitudeQuaternionCov{},
}

var testMpV1Msgs = []Message{
	&MessageHeartbeat{
		Type:           1,
		Autopilot:      2,
		BaseMode:       3,
		CustomMode:     6,
		SystemStatus:   4,
		MavlinkVersion: 5,
	},
	&MessageSysStatus{
		OnboardControlSensorsPresent: 0x01010101,
		OnboardControlSensorsEnabled: 0x01010101,
		OnboardControlSensorsHealth:  0x01010101,
		Load:                         0x0101,
		VoltageBattery:               0x0101,
		CurrentBattery:               0x0101,
		BatteryRemaining:             1,
		DropRateComm:                 0x0101,
		ErrorsComm:                   0x0101,
		ErrorsCount1:                 0x0101,
		ErrorsCount2:                 0x0101,
		ErrorsCount3:                 0x0101,
		ErrorsCount4:                 0x0101,
	},
	&MessageChangeOperatorControl{
		TargetSystem:   1,
		ControlRequest: 1,
		Version:        1,
		Passkey:        "testing",
	},
	&MessageAttitudeQuaternionCov{
		TimeUsec:   2,
		Q:          [4]float32{1, 1, 1, 1},
		Rollspeed:  1,
		Pitchspeed: 1,
		Yawspeed:   1,
		Covariance: [9]float32{1, 1, 1, 1, 1, 1, 1, 1, 1},
	},
}

func TestMessageParserV1Dec(t *testing.T) {
	for i, byt := range testMpV1Bytes {
		mp, _ := newMessageParser(testMpV1Parsers[i])
		msg, err := mp.decode(byt, false)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(msg, testMpV1Msgs[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", msg, testMpV1Msgs[i])
		}
	}
}

func TestMessageParserV1Enc(t *testing.T) {
	for i, msg := range testMpV1Msgs {
		mp, _ := newMessageParser(testMpV1Parsers[i])
		byt, err := mp.encode(msg, false)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(byt, testMpV1Bytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testMpV1Bytes[i])
		}
	}
}

var testMpV2EmptyByteBytes = [][]byte{
	[]byte("\x00\x01\x02\x74\x65\x73\x74\x69\x6e\x67"),
}

var testMpV2EmptyByteParsers = []Message{
	&MessageChangeOperatorControl{},
}

var testMpV2EmptyByteMsgs = []Message{
	&MessageChangeOperatorControl{
		TargetSystem:   0,
		ControlRequest: 1,
		Version:        2,
		Passkey:        "testing",
	},
}

func TestMessageParserV2EmptyByteDec(t *testing.T) {
	for i, byt := range testMpV2EmptyByteBytes {
		mp, _ := newMessageParser(testMpV2EmptyByteParsers[i])
		msg, err := mp.decode(byt, true)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(msg, testMpV2EmptyByteMsgs[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", msg, testMpV2EmptyByteMsgs[i])
		}
	}
}

func TestMessageParserV2EmptyByteEnc(t *testing.T) {
	for i, msg := range testMpV2EmptyByteMsgs {
		mp, _ := newMessageParser(testMpV2EmptyByteParsers[i])
		byt, err := mp.encode(msg, true)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(byt, testMpV2EmptyByteBytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testMpV2EmptyByteBytes[i])
		}
	}
}

var testMpV2ExtensionBytes = [][]byte{
	[]byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x3F\x00\x00\x80\x3F\x00\x00\x80\x3F\x07\x00\x08\x00\x09\x0A\x00\x00\x80\x3F\x00\x00\x80\x3F"),
}

var testMpV2ExtensionParsers = []Message{
	&MessageOpticalFlow{},
}

var testMpV2ExtensionMsgs = []Message{
	&MessageOpticalFlow{
		TimeUsec:       3,
		FlowCompMX:     1,
		FlowCompMY:     1,
		GroundDistance: 1,
		FlowX:          7,
		FlowY:          8,
		SensorId:       9,
		Quality:        0x0A,
		FlowRateX:      1,
		FlowRateY:      1,
	},
}

func TestMessageParserV2ExtensionDec(t *testing.T) {
	for i, byt := range testMpV2ExtensionBytes {
		mp, _ := newMessageParser(testMpV2ExtensionParsers[i])
		msg, err := mp.decode(byt, true)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(msg, testMpV2ExtensionMsgs[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", msg, testMpV2ExtensionMsgs[i])
		}
	}
}

func TestMessageParserV2ExtensionEnc(t *testing.T) {
	for i, msg := range testMpV2ExtensionMsgs {
		mp, _ := newMessageParser(testMpV2ExtensionParsers[i])
		byt, err := mp.encode(msg, true)
		if err != nil {
			t.Fatal(err)
			continue
		}
		if reflect.DeepEqual(byt, testMpV2ExtensionBytes[i]) == false {
			t.Fatalf("invalid: %+v vs %+v", byt, testMpV2ExtensionBytes[i])
		}
	}
}
