package tlog_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v4/pkg/frame"
	"github.com/bluenviron/gomavlib/v4/pkg/message"
	"github.com/bluenviron/gomavlib/v4/pkg/tlog"
)

var casesReadWriter = []struct {
	name string
	dec  []tlog.Entry
	enc  []byte
}{
	{
		"main",
		[]tlog.Entry{
			{
				Time: time.Date(2014, 9, 12, 3, 15, 17, 37000, time.UTC),
				Frame: &frame.V1Frame{
					SequenceNumber: 137,
					SystemID:       17,
					ComponentID:    14,
					Message: &message.MessageRaw{
						ID:      54,
						Payload: []byte{1, 2, 3, 4},
					},
					Checksum: 7854,
				},
			},
			{
				Time: time.Date(2015, 9, 12, 3, 15, 17, 64000, time.UTC),
				Frame: &frame.V2Frame{
					SequenceNumber: 138,
					SystemID:       17,
					ComponentID:    14,
					Message: &message.MessageRaw{
						ID:      54,
						Payload: []byte{1, 2, 3, 4},
					},
					Checksum: 7854,
				},
			},
		},
		[]byte{
			0x00, 0x05, 0x02, 0xd5, 0xb1, 0xc0, 0x1b, 0x65,
			0xfe, 0x04, 0x89, 0x11, 0x0e, 0x36, 0x01, 0x02,
			0x03, 0x04, 0xae, 0x1e, 0x00, 0x05, 0x1f, 0x84,
			0x3d, 0xd3, 0xfb, 0x80, 0xfd, 0x04, 0x00, 0x00,
			0x8a, 0x11, 0x0e, 0x36, 0x00, 0x00, 0x01, 0x02,
			0x03, 0x04, 0xae, 0x1e,
		},
	},
}

func TestReader(t *testing.T) {
	for _, ca := range casesReadWriter {
		t.Run(ca.name, func(t *testing.T) {
			r := &tlog.Reader{
				ByteReader: bytes.NewReader(ca.enc),
			}
			err := r.Initialize()
			require.NoError(t, err)

			i := 0

			for {
				var entry *tlog.Entry
				entry, err = r.Read()
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				require.Equal(t, ca.dec[i], *entry)
				i++
			}

			require.Equal(t, len(ca.dec), i)
		})
	}
}
