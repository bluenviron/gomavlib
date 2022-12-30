package frame

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/message"
)

func TestWriterNewErrors(t *testing.T) {
	_, err := NewWriter(WriterConf{
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.EqualError(t, err, "Writer not provided")

	var buf bytes.Buffer

	_, err = NewWriter(WriterConf{
		Writer:      &buf,
		OutSystemID: 1,
	})
	require.EqualError(t, err, "OutVersion not provided")

	_, err = NewWriter(WriterConf{
		Writer:     &buf,
		OutVersion: V2,
	})
	require.EqualError(t, err, "OutSystemID must be greater than one")

	_, err = NewWriter(WriterConf{
		Writer:     &buf,
		OutVersion: V1,
		OutKey:     NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
	})
	require.EqualError(t, err, "OutSystemID must be greater than one")
}

func TestWriterWriteFrame(t *testing.T) {
	for _, ca := range casesReadWrite {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer, err := NewWriter(WriterConf{
				Writer:      &buf,
				OutVersion:  V2,
				OutSystemID: 1,
				DialectRW:   ca.dialectRW,
			})
			require.NoError(t, err)

			err = writer.WriteFrame(ca.frame)
			require.NoError(t, err)
			require.Equal(t, ca.raw, buf.Bytes())
		})
	}
}

func TestWriterWriteFrameErrors(t *testing.T) {
	for _, ca := range []struct {
		name      string
		dialectRW *dialect.ReadWriter
		frame     Frame
		err       string
	}{
		{
			"nil message",
			nil,
			&V1Frame{
				SequenceID:  0x01,
				SystemID:    0x02,
				ComponentID: 0x03,
				Message:     nil,
				Checksum:    0x0807,
			},
			"message is nil",
		},
		{
			"nil dialect",
			nil,
			&V1Frame{
				SequenceID:  0x27,
				SystemID:    0x01,
				ComponentID: 0x02,
				Message: &MessageTest5{
					'\x10',
					binary.LittleEndian.Uint32([]byte("\x10\x10\x10\x10")),
				},
				Checksum: 0x66e5,
			},
			"dialect is nil",
		},
		{
			"not in dialect",
			testDialectRW,
			&V1Frame{
				SequenceID:  0x27,
				SystemID:    0x01,
				ComponentID: 0x02,
				Message:     &MessageTest8{15, 7},
				Checksum:    0x66e5,
			},
			"message is not in the dialect",
		},
		{
			"frame encode error",
			testDialectRW,
			&V1Frame{
				SequenceID:  0x27,
				SystemID:    0x01,
				ComponentID: 0x02,
				Message:     &MessageTest9{},
				Checksum:    0x66e5,
			},
			"cannot send a message with an ID greater than 255 with a V1 frame",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer, err := NewWriter(WriterConf{
				Writer:      &buf,
				OutVersion:  V2,
				OutSystemID: 1,
				DialectRW:   ca.dialectRW,
			})
			require.NoError(t, err)

			err = writer.WriteFrame(ca.frame)
			require.EqualError(t, err, ca.err)
		})
	}
}

func TestWriterWriteMessage(t *testing.T) {
	// fake current time in order to obtain deterministic signatures
	wayback := time.Date(2019, time.May, 18, 1, 2, 3, 4, time.UTC)
	patch := monkey.Patch(time.Now, func() time.Time { return wayback })
	defer patch.Unpatch()

	for _, c := range []struct {
		name string
		ver  WriterOutVersion
		key  *V2Key
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
			NewV2Key(bytes.Repeat([]byte("\x4F"), 32)),
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
		t.Run(c.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			writer, err := NewWriter(WriterConf{
				Writer:      buf,
				DialectRW:   testDialectRW,
				OutVersion:  c.ver,
				OutSystemID: 1,
				OutKey:      c.key,
			})
			require.NoError(t, err)

			err = writer.WriteMessage(c.msg)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
			buf.Next(len(c.raw))
		})
	}
}

func TestWriterWriteMessageErrors(t *testing.T) {
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
	} {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer, err := NewWriter(WriterConf{
				Writer:      &buf,
				OutVersion:  V2,
				OutSystemID: 1,
				DialectRW:   ca.dialectRW,
			})
			require.NoError(t, err)

			err = writer.WriteMessage(ca.message)
			require.EqualError(t, err, ca.err)
		})
	}
}

func TestWriterWriteFrameNilMsg(t *testing.T) {
	writer, err := NewWriter(WriterConf{
		Writer:      bytes.NewBuffer(nil),
		DialectRW:   nil,
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.NoError(t, err)

	f := &V1Frame{Message: nil}
	err = writer.WriteFrame(f)
	require.Error(t, err)
}
