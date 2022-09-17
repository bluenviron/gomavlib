package frame

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWriterNew(t *testing.T) {
	var buf bytes.Buffer
	_, err := NewReadWriter(ReadWriterConf{
		ReadWriter:  &buf,
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.NoError(t, err)
}
