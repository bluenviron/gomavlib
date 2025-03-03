package streamwriter

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/message"
	"github.com/stretchr/testify/require"
)

var testDialectRW = func() *dialect.ReadWriter {
	d := &dialect.Dialect{
		Version: 3,
		Messages: []message.Message{
			&MessageTest5{},
			&MessageHeartbeat{},
		},
	}

	de := &dialect.ReadWriter{Dialect: d}
	err := de.Initialize()
	if err != nil {
		panic(err)
	}

	return de
}()

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

type MessageTest8 struct {
	TestByte byte
	TestUint uint32
}

func (m *MessageTest8) GetID() uint32 {
	return 8
}

func TestWriteMessage(t *testing.T) {
	// fake current time in order to obtain deterministic signatures
	wayback := time.Date(2019, time.May, 18, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	for _, ca := range []struct {
		name string
		ver  Version
		key  *frame.V2Key
		msg  message.Message
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
			[]byte("\xFD\x09\x01\x00\x00\x01\x01\x00" +
				"\x00\x00\x04\x00\x00\x00\x01\x02" +
				"\x03\x05\x03\x19\xe7\x00\xe0\xf8" +
				"\xd4\xb6\x8e\x0c\xe7\x5d\x07\x46" +
				"\x81\xd4"),
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			fw := &frame.Writer{
				ByteWriter: buf,
				DialectRW:  testDialectRW,
			}
			err := fw.Initialize()
			require.NoError(t, err)

			nw := &Writer{
				FrameWriter: fw,
				Version:     ca.ver,
				SystemID:    1,
				Key:         ca.key,
			}
			err = nw.Initialize()
			require.NoError(t, err)

			err = nw.Write(ca.msg)
			require.NoError(t, err)

			require.Equal(t, ca.raw, buf.Bytes())
		})
	}
}

type MessageTest10 struct {
	TestByte1 byte
	TestUint1 uint32
}

func (m *MessageTest10) GetID() uint32 {
	return 301
}

func TestWriteMessageErrors(t *testing.T) {
	for _, ca := range []struct {
		name      string
		dialectRW *dialect.ReadWriter
		message   message.Message
		err       string
	}{
		{
			"nil message",
			nil,
			nil,
			"message is nil",
		},
		{
			"nil dialect",
			nil,
			&MessageTest8{15, 7},
			"dialect is nil",
		},
		{
			"not in dialect",
			testDialectRW,
			&MessageTest10{15, 7},
			"message is not in the dialect",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)

			fw := &frame.Writer{
				ByteWriter: buf,
				DialectRW:  ca.dialectRW,
			}
			err := fw.Initialize()
			require.NoError(t, err)

			nw := &Writer{
				FrameWriter: fw,
				Version:     V2,
				SystemID:    1,
			}
			err = nw.Initialize()
			require.NoError(t, err)

			err = nw.Write(ca.message)
			require.EqualError(t, err, ca.err)
		})
	}
}
