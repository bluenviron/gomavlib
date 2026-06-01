package frame

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v4/pkg/dialect"
)

func TestWriterNewErrors(t *testing.T) {
	err := (&Writer{}).Initialize()
	require.EqualError(t, err, "ByteWriter not provided")
}

func TestWriterWrite(t *testing.T) {
	for _, ca := range casesReadWrite {
		switch ca.name {
		case "v2 frame with missing empty byte truncation",
			"v1 frame with junk after string termination":
			continue
		}

		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := &Writer{
				ByteWriter: &buf,
				DialectRW:  ca.dialectRW,
			}
			err := writer.Initialize()
			require.NoError(t, err)

			err = writer.Write(ca.frame)
			require.NoError(t, err)
			require.Equal(t, ca.raw, buf.Bytes())
		})
	}
}

func TestWriterWriteErrors(t *testing.T) {
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
				SequenceNumber: 0x01,
				SystemID:       0x02,
				ComponentID:    0x03,
				Message:        nil,
				Checksum:       0x0807,
			},
			"message is nil",
		},
		{
			"nil dialect",
			nil,
			&V1Frame{
				SequenceNumber: 0x27,
				SystemID:       0x01,
				ComponentID:    0x02,
				Message: &MessageTest5{
					'\x10',
					0x10101010,
				},
				Checksum: 0x66e5,
			},
			"dialect is nil",
		},
		{
			"not in dialect",
			testDialectRW,
			&V1Frame{
				SequenceNumber: 0x27,
				SystemID:       0x01,
				ComponentID:    0x02,
				Message:        &MessageTest8{15, 7},
				Checksum:       0x66e5,
			},
			"message is not in the dialect",
		},
		{
			"frame encode error",
			testDialectRW,
			&V1Frame{
				SequenceNumber: 0x27,
				SystemID:       0x01,
				ComponentID:    0x02,
				Message:        &MessageTest9{},
				Checksum:       0x66e5,
			},
			"cannot send a message with an ID greater than 255 with a V1 frame",
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := &Writer{
				ByteWriter: &buf,
				DialectRW:  ca.dialectRW,
			}
			err := writer.Initialize()
			require.NoError(t, err)

			err = writer.Write(ca.frame)
			require.EqualError(t, err, ca.err)
		})
	}
}

func TestWriterWriteFrameNilMsg(t *testing.T) {
	writer := &Writer{
		ByteWriter: bytes.NewBuffer(nil),
	}
	err := writer.Initialize()
	require.NoError(t, err)

	f := &V1Frame{Message: nil}
	err = writer.Write(f)
	require.Error(t, err)
}
