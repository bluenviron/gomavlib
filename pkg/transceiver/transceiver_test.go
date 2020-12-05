package transceiver

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/frame"
	"github.com/aler9/gomavlib/pkg/msg"
)

type MAV_TYPE int      //nolint:golint
type MAV_AUTOPILOT int //nolint:golint
type MAV_MODE_FLAG int //nolint:golint
type MAV_STATE int     //nolint:golint

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
	SensorId       uint8 //nolint:golint
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

var testDialectDE = func() *dialect.DecEncoder {
	d := &dialect.Dialect{3, []msg.Message{ //nolint:govet
		&MessageTest5{},
		&MessageTest6{},
		&MessageTest8{},
		&MessageHeartbeat{},
		&MessageOpticalFlow{},
	}}
	de, err := dialect.NewDecEncoder(d)
	if err != nil {
		panic(err)
	}
	return de
}()

var casesTransceiver = []struct {
	name      string
	dialectDE *dialect.DecEncoder
	key       *frame.V2Key
	frame     frame.Frame
	raw       []byte
}{
	{
		"v1 frame with nil content",
		nil,
		nil,
		&frame.V1Frame{
			SequenceID:  0x01,
			SystemID:    0x02,
			ComponentID: 0x03,
			Message: &msg.MessageRaw{
				ID:      4,
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
		&frame.V1Frame{
			SequenceID:  0x27,
			SystemID:    0x01,
			ComponentID: 0x02,
			Message: &msg.MessageRaw{ //nolint:govet
				8,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0xc7fa,
		},
		[]byte("\xFE\x05\x27\x01\x02\x08\x10\x10\x10\x10\x10\xfa\xc7"),
	},
	{
		"v1 frame with decoded message",
		testDialectDE,
		nil,
		&frame.V1Frame{
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
		testDialectDE,
		nil,
		&frame.V2Frame{
			IncompatibilityFlag: 0,
			CompatibilityFlag:   0,
			SequenceID:          3,
			SystemID:            4,
			ComponentID:         5,
			Message: &msg.MessageRaw{ //nolint:govet
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
		&frame.V2Frame{
			IncompatibilityFlag: 0x00,
			CompatibilityFlag:   0x00,
			SequenceID:          0x8F,
			SystemID:            0x01,
			ComponentID:         0x02,
			Message: &msg.MessageRaw{ //nolint:govet
				0x0607,
				[]byte("\x10\x10\x10\x10\x10"),
			},
			Checksum: 0x0349,
		},
		[]byte("\xFD\x05\x00\x00\x8F\x01\x02\x07\x06\x00\x10\x10\x10\x10\x10\x49\x03"),
	},
	{
		"v2 frame with decoded message",
		testDialectDE,
		nil,
		&frame.V2Frame{
			IncompatibilityFlag: 0x00,
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
		testDialectDE,
		frame.NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
		&frame.V2Frame{
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
			Signature:          &frame.V2Signature{0x0e, 0x47, 0x04, 0x0c, 0xef, 0x9b},
		},
		[]byte("\xFD\x09\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0e\x47\x04\x0c\xef\x9b"),
	},
	{
		"v2 frame with decoded message, signed",
		testDialectDE,
		frame.NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
		&frame.V2Frame{
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
			Signature:          &frame.V2Signature{0xa8, 0x88, 0x9, 0x39, 0xb2, 0x60},
		},
		[]byte("\xFD\x22\x01\x00\x00\x00\x00\x64\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\xa0\x40\x00\x00\xc0\x40\x00\x00\x00\x41\x03\x00\x04\x00\x02\x07\x00\x00\x00\x00\x00\x00\x80\x3f\x77\xfb\x03\x04\x00\x00\x00\x00\x00\xa8\x88\x09\x39\xb2\x60"),
	},
}

func TestTransceiverDecode(t *testing.T) {
	for _, c := range casesTransceiver {
		t.Run(c.name, func(t *testing.T) {
			transceiver, err := New(Conf{
				Reader:      bytes.NewReader(c.raw),
				Writer:      bytes.NewBuffer(nil),
				DialectDE:   c.dialectDE,
				OutVersion:  V2,
				OutSystemID: 1,
				InKey:       c.key,
			})
			require.NoError(t, err)
			frame, err := transceiver.Read()
			require.NoError(t, err)
			require.Equal(t, c.frame, frame)
		})
	}
}

func TestTransceiverEncode(t *testing.T) {
	for _, c := range casesTransceiver {
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			transceiver, err := New(Conf{
				Reader:      bytes.NewBuffer(nil),
				Writer:      buf,
				OutVersion:  V2,
				OutSystemID: 1,
				DialectDE:   c.dialectDE,
			})
			require.NoError(t, err)
			err = transceiver.WriteFrame(c.frame)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
		})
	}
}

var casesTransceiverWriteMessage = []struct {
	name string
	ver  Version
	key  *frame.V2Key
	msg  msg.Message
	raw  []byte
}{
	{
		"v1 frame",
		V1,
		nil,
		&MessageTest5{
			'\x10',
			binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
		},
		[]byte("\xFE\x05\x00\x01\x01\x05\x10\x10\x10\x10\x10\x75\x84"),
	},
	{
		"v2 frame, signed",
		V2,
		frame.NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
		&MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     4,
			SystemStatus:   5,
			MavlinkVersion: 3,
		},
		[]byte("\xFD\x09\x01\x00\x00\x01\x01\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\x19\xe7\x00\xe0\xf8\xd4\xb6\x8e\x0c\xe7\x5d\x07\x46\x81\xd4"),
	},
}

func TestTransceiverWriteMessage(t *testing.T) {
	// fake current time in order to obtain deterministic signatures
	wayback := time.Date(2019, time.May, 18, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	for _, c := range casesTransceiverWriteMessage {
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			transceiver, err := New(Conf{
				Reader:      bytes.NewBuffer(nil),
				Writer:      buf,
				DialectDE:   testDialectDE,
				OutVersion:  c.ver,
				OutSystemID: 1,
				OutKey:      c.key,
			})
			require.NoError(t, err)

			err = transceiver.WriteMessage(c.msg)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
			buf.Next(len(c.raw))
		})
	}
}

func TestTransceiverEncodeNilMsg(t *testing.T) {
	transceiver, err := New(Conf{
		Reader:      bytes.NewReader(nil),
		Writer:      bytes.NewBuffer(nil),
		DialectDE:   nil,
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.NoError(t, err)

	f := &frame.V1Frame{Message: nil}
	err = transceiver.WriteFrame(f)
	require.Error(t, err)
}

// ensure that the Frame is left untouched by WriteFrame()
// and therefore the function can be called by multiple routines in parallel
func TestTransceiverWriteFrameIsConst(t *testing.T) {
	dialectDE, err := dialect.NewDecEncoder(&dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}}) //nolint:govet
	require.NoError(t, err)

	transceiver, err := New(Conf{
		Reader:      bytes.NewReader(nil),
		Writer:      bytes.NewBuffer(nil),
		DialectDE:   dialectDE,
		OutVersion:  V2,
		OutSystemID: 1,
		OutKey:      frame.NewV2Key(bytes.Repeat([]byte("\x7C"), 32)),
	})
	require.NoError(t, err)

	f := &frame.V2Frame{
		Message: &MessageHeartbeat{
			Type:           1,
			Autopilot:      2,
			BaseMode:       3,
			CustomMode:     4,
			SystemStatus:   5,
			MavlinkVersion: 3,
		},
	}
	original := f.Clone()

	err = transceiver.WriteFrame(f)
	require.NoError(t, err)
	require.Equal(t, f, original)
}
