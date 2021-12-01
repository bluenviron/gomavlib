package parser

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

func TestWriterWriteFrame(t *testing.T) {
	for _, c := range casesReadWrite {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer, err := NewWriter(WriterConf{
				Writer:      &buf,
				OutVersion:  V2,
				OutSystemID: 1,
				DialectDE:   c.dialectDE,
			})
			require.NoError(t, err)
			err = writer.WriteFrame(c.frame)
			require.NoError(t, err)
			require.Equal(t, c.raw, buf.Bytes())
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
				DialectDE:   testDialectDE,
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

func TestWriterWriteFrameNilMsg(t *testing.T) {
	writer, err := NewWriter(WriterConf{
		Writer:      bytes.NewBuffer(nil),
		DialectDE:   nil,
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.NoError(t, err)

	f := &frame.V1Frame{Message: nil}
	err = writer.WriteFrame(f)
	require.Error(t, err)
}

// ensure that the Frame is left untouched by WriteFrame()
// and therefore the function can be called by multiple routines in parallel
func TestWriterWriteFrameIsConst(t *testing.T) {
	dialectDE, err := dialect.NewDecEncoder(&dialect.Dialect{3, []msg.Message{&MessageHeartbeat{}}}) //nolint:govet
	require.NoError(t, err)

	writer, err := NewWriter(WriterConf{
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

	err = writer.WriteFrame(f)
	require.NoError(t, err)
	require.Equal(t, f, original)
}
