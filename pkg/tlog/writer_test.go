package tlog_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v4/pkg/tlog"
)

func TestWriter(t *testing.T) {
	for _, ca := range casesReadWriter {
		t.Run(ca.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := &tlog.Writer{ByteWriter: &buf}
			err := w.Initialize()
			require.NoError(t, err)

			for _, entry := range ca.dec {
				err = w.Write(&entry)
				require.NoError(t, err)
			}

			require.Equal(t, ca.enc, buf.Bytes())
		})
	}
}
