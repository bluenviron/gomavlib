package frame

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWriter(t *testing.T) {
	var buf bytes.Buffer
	_, err := NewReadWriter(ReadWriterConf{
		ReadWriter:  &buf,
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.NoError(t, err)
}

func TestReadWriterErrors(t *testing.T) {
	_, err := NewReadWriter(ReadWriterConf{
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.EqualError(t, err, "BufByteReader not provided")

	var buf bytes.Buffer
	_, err = NewReadWriter(ReadWriterConf{
		ReadWriter:  &buf,
		OutSystemID: 1,
	})
	require.EqualError(t, err, "OutVersion not provided")
}
