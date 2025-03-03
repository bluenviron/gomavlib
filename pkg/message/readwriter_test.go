package message_test

import (
	"bytes"
	"testing"

	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
	"github.com/stretchr/testify/require"
)

type (
	MAV_TYPE              uint64 //nolint:revive
	MAV_AUTOPILOT         uint64 //nolint:revive
	MAV_MODE_FLAG         uint64 //nolint:revive
	MAV_STATE             uint64 //nolint:revive
	MAV_SYS_STATUS_SENSOR uint64 //nolint:revive
	MAV_CMD               uint64 //nolint:revive
	MYENUM                uint64
)

type MessageAllTypes struct {
	A uint8
	B int8
	C uint16
	D int16
	E uint32
	F int32
	G uint64
	H int64
	I float32
	J float64
	K string `mavlen:"30"`
	L MYENUM `mavenum:"uint8"`
	M MYENUM `mavenum:"int8"`
	N MYENUM `mavenum:"uint16"`
	P MYENUM `mavenum:"uint32"`
	Q MYENUM `mavenum:"int32"`
	R MYENUM `mavenum:"uint64"`
	S string
}

func (*MessageAllTypes) GetID() uint32 {
	return 155
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

type MessageSysStatus struct {
	OnboardControlSensorsPresent MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
	OnboardControlSensorsEnabled MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
	OnboardControlSensorsHealth  MAV_SYS_STATUS_SENSOR `mavenum:"uint32"`
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

func (m *MessageSysStatus) GetID() uint32 {
	return 1
}

type MessageChangeOperatorControl struct {
	TargetSystem   uint8
	ControlRequest uint8
	Version        uint8
	Passkey        string `mavlen:"25"`
}

func (m *MessageChangeOperatorControl) GetID() uint32 {
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

func (m *MessageAttitudeQuaternionCov) GetID() uint32 {
	return 61
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

type MessagePlayTune struct {
	TargetSystem    uint8
	TargetComponent uint8
	Tune            string `mavlen:"30"`
	Tune2           string `mavext:"true" mavlen:"200"`
}

func (*MessagePlayTune) GetID() uint32 {
	return 258
}

type MessageAhrs struct {
	OmegaIx     float32 `mavname:"omegaIx"`
	OmegaIy     float32 `mavname:"omegaIy"`
	OmegaIz     float32 `mavname:"omegaIz"`
	AccelWeight float32
	RenormVal   float32
	ErrorRp     float32
	ErrorYaw    float32
}

func (*MessageAhrs) GetID() uint32 {
	return 163
}

type MessageTrajectoryRepresentationWaypoints struct {
	TimeUsec    uint64
	ValidPoints uint8
	PosX        [5]float32
	PosY        [5]float32
	PosZ        [5]float32
	VelX        [5]float32
	VelY        [5]float32
	VelZ        [5]float32
	AccX        [5]float32
	AccY        [5]float32
	AccZ        [5]float32
	PosYaw      [5]float32
	VelYaw      [5]float32
	Command     [5]MAV_CMD `mavenum:"uint16"`
}

func (*MessageTrajectoryRepresentationWaypoints) GetID() uint32 {
	return 332
}

var casesCRC = []struct {
	msg message.Message
	crc byte
}{
	{
		&MessageHeartbeat{},
		50,
	},
	{
		&MessageSysStatus{},
		124,
	},
	{
		&MessageChangeOperatorControl{},
		217,
	},
	{
		&MessageAttitudeQuaternionCov{},
		167,
	},
	{
		&MessageOpticalFlow{},
		175,
	},
	{
		&MessagePlayTune{},
		187,
	},
	{
		&MessageAhrs{},
		127,
	},
}

func TestCRC(t *testing.T) {
	for _, c := range casesCRC {
		mp, err := message.NewReadWriter(c.msg)
		require.NoError(t, err)
		require.Equal(t, c.crc, mp.CRCExtra())
	}
}

var casesReadWriter = []struct {
	name   string
	isV2   bool
	parsed message.Message
	raw    []byte
}{
	{
		"v1",
		false,
		&MessageAllTypes{
			A: 127,
			B: -12,
			C: 1343,
			D: 5652,
			E: 5323,
			F: 7987,
			G: 8654,
			H: 6753,
			I: 5764,
			J: 3423,
			K: "teststring",
			L: 232,
			M: 34,
			N: 1422,
			P: 1233,
			Q: 2343,
			R: 1232,
			S: "a",
		},
		[]byte{
			0xce, 0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x61, 0x1a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0xbe, 0xaa, 0x40,
			0xd0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0xcb, 0x14, 0x0, 0x0, 0x33, 0x1f, 0x0, 0x0,
			0x0, 0x20, 0xb4, 0x45, 0xd1, 0x4, 0x0, 0x0,
			0x27, 0x9, 0x0, 0x0, 0x3f, 0x5, 0x14, 0x16,
			0x8e, 0x5, 0x7f, 0xf4, 0x74, 0x65, 0x73, 0x74,
			0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x0, 0x0, 0xe8, 0x22, 0x61,
		},
	},
	{
		"v1 with string with max length",
		false,
		&MessageChangeOperatorControl{
			Passkey: "abcdefghijklmnopqrstuvwxy",
		},
		[]byte("\x00\x00\x00\x61\x62\x63\x64\x65" +
			"\x66\x67\x68\x69\x6a\x6b\x6c\x6d" +
			"\x6e\x6f\x70\x71\x72\x73\x74\x75" +
			"\x76\x77\x78\x79"),
	},
	{
		"v1 with array",
		false,
		&MessageAttitudeQuaternionCov{
			TimeUsec:   2,
			Q:          [4]float32{1, 1, 1, 1},
			Rollspeed:  1,
			Pitchspeed: 1,
			Yawspeed:   1,
			Covariance: [9]float32{1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		append(
			[]byte("\x02\x00\x00\x00\x00\x00\x00\x00"),
			bytes.Repeat([]byte("\x00\x00\x80\x3F"), 16)...),
	},
	{
		"v1 with extensions",
		false,
		&MessageOpticalFlow{
			TimeUsec:       3,
			FlowCompMX:     1,
			FlowCompMY:     1,
			GroundDistance: 1,
			FlowX:          7,
			FlowY:          8,
			SensorId:       9,
			Quality:        0x0A,
		},
		[]byte("\x03\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x80\x3F\x00\x00\x80\x3F" +
			"\x00\x00\x80\x3F\x07\x00\x08\x00" +
			"\x09\x0A"),
	},
	{
		"v1 with array of enums",
		false,
		&MessageTrajectoryRepresentationWaypoints{
			TimeUsec:    1,
			ValidPoints: 2,
			PosX:        [5]float32{1, 2, 3, 4, 5},
			PosY:        [5]float32{1, 2, 3, 4, 5},
			PosZ:        [5]float32{1, 2, 3, 4, 5},
			VelX:        [5]float32{1, 2, 3, 4, 5},
			VelY:        [5]float32{1, 2, 3, 4, 5},
			VelZ:        [5]float32{1, 2, 3, 4, 5},
			AccX:        [5]float32{1, 2, 3, 4, 5},
			AccY:        [5]float32{1, 2, 3, 4, 5},
			AccZ:        [5]float32{1, 2, 3, 4, 5},
			PosYaw:      [5]float32{1, 2, 3, 4, 5},
			VelYaw:      [5]float32{1, 2, 3, 4, 5},
			Command:     [5]MAV_CMD{1, 2, 3, 4, 5},
		},
		[]byte("\x01\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x00\x00\x80\x3f" +
			"\x00\x00\x00\x40\x00\x00\x40\x40" +
			"\x00\x00\x80\x40\x00\x00\xa0\x40" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x00\x00\x80\x3f" +
			"\x00\x00\x00\x40\x00\x00\x40\x40" +
			"\x00\x00\x80\x40\x00\x00\xa0\x40" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x00\x00\x80\x3f" +
			"\x00\x00\x00\x40\x00\x00\x40\x40" +
			"\x00\x00\x80\x40\x00\x00\xa0\x40" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x00\x00\x80\x3f" +
			"\x00\x00\x00\x40\x00\x00\x40\x40" +
			"\x00\x00\x80\x40\x00\x00\xa0\x40" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x00\x00\x80\x3f" +
			"\x00\x00\x00\x40\x00\x00\x40\x40" +
			"\x00\x00\x80\x40\x00\x00\xa0\x40" +
			"\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40\x01\x00\x02\x00" +
			"\x03\x00\x04\x00\x05\x00\x02"),
	},
	{
		"v2 with empty-byte truncation and empty",
		true,
		&MessageAhrs{},
		[]byte("\x00"),
	},
	{
		"v2 with empty-byte truncation a",
		true,
		&MessageChangeOperatorControl{
			TargetSystem:   0,
			ControlRequest: 1,
			Version:        2,
			Passkey:        "testing",
		},
		[]byte("\x00\x01\x02\x74\x65\x73\x74\x69" +
			"\x6e\x67"),
	},
	{
		"v2 with empty-byte truncation b",
		true,
		&MessageAhrs{
			OmegaIx:     1,
			OmegaIy:     2,
			OmegaIz:     3,
			AccelWeight: 4,
			RenormVal:   5,
			ErrorRp:     0,
			ErrorYaw:    0,
		},
		[]byte("\x00\x00\x80\x3f\x00\x00\x00\x40" +
			"\x00\x00\x40\x40\x00\x00\x80\x40" +
			"\x00\x00\xa0\x40"),
	},
	{
		"v2 with extensions a",
		true,
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
		[]byte("\x03\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x80\x3F\x00\x00\x80\x3F" +
			"\x00\x00\x80\x3F\x07\x00\x08\x00" +
			"\x09\x0A\x00\x00\x80\x3F\x00\x00" +
			"\x80\x3F"),
	},
	{
		"v2 with extensions b",
		true,
		&MessagePlayTune{
			TargetSystem:    1,
			TargetComponent: 2,
			Tune:            "test1",
			Tune2:           "test2",
		},
		[]byte("\x01\x02\x74\x65\x73\x74\x31\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x00\x00\x00\x00\x00\x00\x00\x00" +
			"\x74\x65\x73\x74\x32"),
	},
}

type Invalid struct{}

func (*Invalid) GetID() uint32 {
	return 0
}

type MYENUM2 int8

type MessageInvalidEnum struct {
	MyEnum MYENUM2 `mavenum:"int8"`
}

func (*MessageInvalidEnum) GetID() uint32 {
	return 0
}

type MYENUM3 uint64

type MessageInvalidEnum2 struct {
	MyEnum MYENUM3 `mavenum:"invalid"`
}

func (*MessageInvalidEnum2) GetID() uint32 {
	return 0
}

type MYENUM4 uint64

type MessageInvalidEnum3 struct {
	MyEnum MYENUM4 `mavenum:"int64"`
}

func (*MessageInvalidEnum3) GetID() uint32 {
	return 0
}

type MessageInvalid2 struct {
	Pointer *int
}

func (*MessageInvalid2) GetID() uint32 {
	return 0
}

type MessageInvalid3 struct {
	Str string `mavlen:"invalid"`
}

func (*MessageInvalid3) GetID() uint32 {
	return 0
}

func TestNewReadWriterErrors(t *testing.T) {
	_, err := message.NewReadWriter(&Invalid{})
	require.EqualError(t, err, "struct name must begin with 'Message'")

	_, err = message.NewReadWriter(&MessageInvalidEnum{})
	require.EqualError(t, err, "an enum must be an uint64")

	_, err = message.NewReadWriter(&MessageInvalidEnum2{})
	require.EqualError(t, err, "unsupported Go type: invalid")

	_, err = message.NewReadWriter(&MessageInvalidEnum3{})
	require.EqualError(t, err, "type 'int64' cannot be used as enum")

	_, err = message.NewReadWriter(&MessageInvalid2{})
	require.EqualError(t, err, "unsupported Go type: ")

	_, err = message.NewReadWriter(&MessageInvalid3{})
	require.EqualError(t, err, "string has invalid length: invalid")
}

func TestRead(t *testing.T) {
	for _, c := range casesReadWriter {
		t.Run(c.name, func(t *testing.T) {
			mp, err := message.NewReadWriter(c.parsed)
			require.NoError(t, err)
			msg, err := mp.Read(&message.MessageRaw{
				ID:      c.parsed.GetID(),
				Payload: c.raw,
			}, c.isV2)
			require.NoError(t, err)
			require.Equal(t, c.parsed, msg)
		})
	}
}

func TestWrite(t *testing.T) {
	for _, c := range casesReadWriter {
		t.Run(c.name, func(t *testing.T) {
			mp, err := message.NewReadWriter(c.parsed)
			require.NoError(t, err)
			msgRaw := mp.Write(c.parsed, c.isV2)
			require.Equal(t, &message.MessageRaw{
				ID:      c.parsed.GetID(),
				Payload: c.raw,
			}, msgRaw)
		})
	}
}

func FuzzReadWriter(f *testing.F) {
	for _, ca := range casesReadWriter {
		f.Add(ca.raw, ca.parsed.GetID(), false, false)
	}

	f.Fuzz(func(t *testing.T, raw []byte, msgID uint32, v2In bool, v2Out bool) {
		if msgID >= uint32(len(ardupilotmega.Dialect.Messages)) {
			return
		}

		rw := &message.ReadWriter{Message: ardupilotmega.Dialect.Messages[msgID]}
		err := rw.Initialize()
		require.NoError(t, err)

		msg, err := rw.Read(&message.MessageRaw{
			ID:      msgID,
			Payload: raw,
		}, v2In)
		if err == nil {
			rw.Write(msg, v2Out)
		}
	})
}
