package gomavlib

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var casesCRC = []struct {
	msg Message
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

func TestDialectCRC(t *testing.T) {
	for _, c := range casesCRC {
		mp, err := newDialectMessage(c.msg)
		require.NoError(t, err)
		require.Equal(t, c.crc, mp.crcExtra)
	}
}

var casesMsgs = []struct {
	name   string
	isV2   bool
	parsed Message
	raw    []byte
}{
	{
		"v1 message",
		false,
		&MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     6,
			SystemStatus:   4,
			MavlinkVersion: 5,
		},
		[]byte("\x06\x00\x00\x00\x01\x02\x03\x04\x05"),
	},
	{
		"v1 message",
		false,
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
		bytes.Repeat([]byte("\x01"), 31),
	},
	{
		"v1 message",
		false,
		&MessageChangeOperatorControl{
			TargetSystem:   1,
			ControlRequest: 1,
			Version:        1,
			Passkey:        "testing",
		},
		[]byte("\x01\x01\x01\x74\x65\x73\x74\x69\x6e\x67\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
	},
	{
		"v1 message with array",
		false,
		&MessageAttitudeQuaternionCov{
			TimeUsec:   2,
			Q:          [4]float32{1, 1, 1, 1},
			Rollspeed:  1,
			Pitchspeed: 1,
			Yawspeed:   1,
			Covariance: [9]float32{1, 1, 1, 1, 1, 1, 1, 1, 1},
		},
		append([]byte("\x02\x00\x00\x00\x00\x00\x00\x00"), bytes.Repeat([]byte("\x00\x00\x80\x3F"), 16)...),
	},
	{
		"v1 message with extension fields (skipped)",
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
		[]byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x3F\x00\x00\x80\x3F\x00\x00\x80\x3F\x07\x00\x08\x00\x09\x0A"),
	},
	{
		"v1 message with array of enums",
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
		[]byte("\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40\x01\x00\x02\x00\x03\x00\x04\x00\x05\x00\x02"),
	},
	{
		"v2 message with empty-byte truncation and empty",
		true,
		&MessageAhrs{},
		[]byte("\x00"),
	},
	{
		"v2 message with empty-byte truncation",
		true,
		&MessageChangeOperatorControl{
			TargetSystem:   0,
			ControlRequest: 1,
			Version:        2,
			Passkey:        "testing",
		},
		[]byte("\x00\x01\x02\x74\x65\x73\x74\x69\x6e\x67"),
	},
	{
		"v2 message with empty-byte truncation",
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
		[]byte("\x00\x00\x80\x3f\x00\x00\x00\x40\x00\x00\x40\x40\x00\x00\x80\x40\x00\x00\xa0\x40"),
	},
	{
		"v2 message with extensions",
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
		[]byte("\x03\x00\x00\x00\x00\x00\x00\x00\x00\x00\x80\x3F\x00\x00\x80\x3F\x00\x00\x80\x3F\x07\x00\x08\x00\x09\x0A\x00\x00\x80\x3F\x00\x00\x80\x3F"),
	},
	{
		"v2 message with extensions",
		true,
		&MessagePlayTune{
			TargetSystem:    1,
			TargetComponent: 2,
			Tune:            "test1",
			Tune2:           "test2",
		},
		[]byte("\x01\x02\x74\x65\x73\x74\x31\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x74\x65\x73\x74\x32"),
	},
}

func TestDialectDecode(t *testing.T) {
	for _, c := range casesMsgs {
		t.Run(c.name, func(t *testing.T) {
			mp, err := newDialectMessage(c.parsed)
			require.NoError(t, err)
			msg, err := mp.decode(c.raw, c.isV2)
			require.NoError(t, err)
			require.Equal(t, c.parsed, msg)
		})
	}
}

func TestDialectEncode(t *testing.T) {
	for _, c := range casesMsgs {
		t.Run(c.name, func(t *testing.T) {
			mp, err := newDialectMessage(c.parsed)
			require.NoError(t, err)
			byt, err := mp.encode(c.parsed, c.isV2)
			require.NoError(t, err)
			require.Equal(t, c.raw, byt)
		})
	}
}
