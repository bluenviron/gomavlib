package frame

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadWriterNew(t *testing.T) {
	var buf bytes.Buffer
	_, err := NewReadWriter(
		ReaderConf{
			Reader: &buf,
		},
		WriterConf{
			Writer:      &buf,
			OutVersion:  V2,
			OutSystemID: 1,
		})
	require.NoError(t, err)
}
