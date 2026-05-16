package frame

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
)

type (
	MAV_TYPE      uint64 //nolint:revive
	MAV_AUTOPILOT uint64 //nolint:revive
	MAV_MODE_FLAG uint64 //nolint:revive
	MAV_STATE     uint64 //nolint:revive
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

type MAV_CMD uint64 //nolint:revive

type MessageCommandLong struct {
	// System which should execute the command
	TargetSystem uint8
	// Component which should execute the command, 0 for all components
	TargetComponent uint8
	// Command ID (of command to send).
	Command MAV_CMD `mavenum:"uint16"`
	// 0: First transmission of this command. 1-255: Confirmation transmissions (e.g. for kill command)
	Confirmation uint8
	// Parameter 1 (for the specific command).
	Param1 float32
	// Parameter 2 (for the specific command).
	Param2 float32
	// Parameter 3 (for the specific command).
	Param3 float32
	// Parameter 4 (for the specific command).
	Param4 float32
	// Parameter 5 (for the specific command).
	Param5 float32
	// Parameter 6 (for the specific command).
	Param6 float32
	// Parameter 7 (for the specific command).
	Param7 float32
}

// GetID implements the message.Message interface.
func (*MessageCommandLong) GetID() uint32 {
	return 76
}

type MAV_PARAM_TYPE uint64 //nolint:revive

const (
	MAV_PARAM_TYPE_INT16 MAV_PARAM_TYPE = 4 //nolint:revive
)

type MessageParamSet struct {
	TargetSystem    uint8
	TargetComponent uint8
	ParamId         string `mavlen:"16"` //nolint:revive
	ParamValue      float32
	ParamType       MAV_PARAM_TYPE `mavenum:"uint8"`
}

// GetID implements the message.Message interface.
func (*MessageParamSet) GetID() uint32 {
	return 23
}

var testDialectRW = func() *dialect.ReadWriter {
	d := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageTest5{},
			&MessageTest6{},
			&MessageTest9{},
			&MessageHeartbeat{},
			&MessageOpticalFlow{},
		},
	}

	de := &dialect.ReadWriter{Dialect: d}
	err := de.Initialize()
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
			SequenceNumber: 0x01,
			SystemID:       0x02,
			ComponentID:    0x03,
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
			SequenceNumber: 0x27,
			SystemID:       0x01,
			ComponentID:    0x02,
			Message: &message.MessageRaw{
				ID:      8,
				Payload: []byte("\x10\x10\x10\x10\x10"),
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
			SequenceNumber: 0x27,
			SystemID:       0x01,
			ComponentID:    0x02,
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
			SequenceNumber:      3,
			SystemID:            4,
			ComponentID:         5,
			Message: &message.MessageRaw{
				ID:      4,
				Payload: nil,
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
			SequenceNumber:      0x8F,
			SystemID:            0x01,
			ComponentID:         0x02,
			Message: &message.MessageRaw{
				ID:      0x0607,
				Payload: []byte("\x10\x10\x10\x10\x10"),
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
			SequenceNumber:      0x8F,
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
			SequenceNumber:      0x00,
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
			SequenceNumber:      0x00,
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
	{
		"v2 frame with missing empty byte truncation",
		func() *dialect.ReadWriter {
			d := &dialect.Dialect{
				Version: 3,
				Messages: []message.Message{
					&MessageCommandLong{},
				},
			}

			drw := &dialect.ReadWriter{Dialect: d}
			err := drw.Initialize()
			if err != nil {
				panic(err)
			}

			return drw
		}(),
		nil,
		&V2Frame{
			SequenceNumber: 11,
			SystemID:       255,
			ComponentID:    220,
			Message: &MessageCommandLong{
				TargetSystem:    1,
				TargetComponent: 1,
				Command:         176,
				Param1:          1,
			},
			Checksum: 51683,
		},
		[]byte{
			0xfd, 0x21, 0x00, 0x00, 0x0b, 0xff, 0xdc, 0x4c,
			0x00, 0x00, 0x00, 0x00, 0x80, 0x3f, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xb0, 0x00,
			0x01, 0x01, 0x00, 0xc4, 0x75,
		},
	},
	{
		"v1 frame with junk after string termination",
		func() *dialect.ReadWriter {
			d := &dialect.Dialect{
				Version: 3,
				Messages: []message.Message{
					&MessageParamSet{},
				},
			}

			drw := &dialect.ReadWriter{Dialect: d}
			err := drw.Initialize()
			if err != nil {
				panic(err)
			}

			return drw
		}(),
		nil,
		&V1Frame{
			SequenceNumber: 11,
			SystemID:       255,
			ComponentID:    220,
			Message: &MessageParamSet{
				TargetSystem:    12,
				TargetComponent: 13,
				ParamId:         "RTL_ALT",
				ParamValue:      4936,
				ParamType:       MAV_PARAM_TYPE_INT16,
			},
			Checksum: 0x745c,
		},
		[]byte{
			0xfe, 0x17, 0x0b, 0xff, 0xdc, 0x17, 0x00, 0x40,
			0x9a, 0x45, 0x0c, 0x0d, 0x52, 0x54, 0x4c, 0x5f,
			0x41, 0x4c, 0x54, 0x00, 0x01, 0x02, 0x03, 0x00,
			0x00, 0x00, 0x00, 0x00, 0x04, 0x51, 0x49,
		},
	},
}

func TestReaderNewErrors(t *testing.T) {
	_, err := NewReader(ReaderConf{})
	require.EqualError(t, err, "BufByteReader not provided")
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

func TestReaderErrorSignatureTimestamp(t *testing.T) {
	var buf bytes.Buffer

	msgByts := []byte{4, 0, 0, 0, 1, 2, 3, 5, 3}

	f := &V2Frame{
		IncompatibilityFlag: 0x01,
		CompatibilityFlag:   0x00,
		SequenceNumber:      0x00,
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
	n, err := f.marshalTo(buf2, msgByts)
	require.NoError(t, err)
	buf.Write(buf2[:n])

	f = &V2Frame{
		IncompatibilityFlag: 0x01,
		CompatibilityFlag:   0x00,
		SequenceNumber:      0x00,
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
	n, err = f.marshalTo(buf2, msgByts)
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

// datagramReader mimics the semantics of a UDP net.Conn: each Read() returns
// exactly one datagram and any leftover bytes that did not fit into the caller's
// buffer are silently discarded by the kernel, as documented for SOCK_DGRAM.
type datagramReader struct {
	datagrams [][]byte
}

func (d *datagramReader) Read(p []byte) (int, error) {
	if len(d.datagrams) == 0 {
		return 0, io.EOF
	}
	dg := d.datagrams[0]
	d.datagrams = d.datagrams[1:]
	n := copy(p, dg)
	// emulate the kernel discarding the tail of an oversized datagram.
	return n, nil
}

// TestReaderUDPMultiFrameDatagram verifies that the Reader can decode every
// frame inside a single UDP datagram even when their combined size exceeds the
// classic MAVLink frame budget. Routers and autopilots (e.g. ArduPilot via
// mavproxy / mavp2p forwarders) routinely coalesce multiple frames into one
// datagram, and an undersized internal buffer would cause the kernel to
// truncate the datagram, producing spurious "invalid magic byte" and
// "wrong checksum" errors. Regression test for the UDP parsing failure
// reported against valid ArduPilot streams.
func TestReaderUDPMultiFrameDatagram(t *testing.T) {
	// craft three frames whose total encoded length is larger than the
	// previous internal buffer size (512 bytes).
	makeFrame := func(seq byte, payloadLen int) ([]byte, *V2Frame) {
		payload := bytes.Repeat([]byte{0x10}, payloadLen)
		f := &V2Frame{
			IncompatibilityFlag: 0,
			CompatibilityFlag:   0,
			SequenceNumber:      seq,
			SystemID:            1,
			ComponentID:         2,
			Message: &message.MessageRaw{
				ID:      0x0607,
				Payload: payload,
			},
		}
		f.Checksum = f.GenerateChecksum(0)
		buf := make([]byte, 512)
		n, err := f.marshalTo(buf, payload)
		require.NoError(t, err)
		return buf[:n], f
	}

	raw1, want1 := makeFrame(1, 250)
	raw2, want2 := makeFrame(2, 250)
	raw3, want3 := makeFrame(3, 50)

	// concatenate every frame into a single UDP datagram (>512 bytes).
	datagram := append(append(append([]byte{}, raw1...), raw2...), raw3...)
	require.Greater(t, len(datagram), 512)

	reader, err := NewReader(ReaderConf{
		Reader: &datagramReader{datagrams: [][]byte{datagram}},
	})
	require.NoError(t, err)

	for _, want := range []*V2Frame{want1, want2, want3} {
		got, errR := reader.Read()
		require.NoError(t, errR)
		require.Equal(t, want, got)
	}
}
