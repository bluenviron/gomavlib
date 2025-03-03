package frame

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriterOutVersion(t *testing.T) {
	require.NotNil(t, V1.String())
	require.NotNil(t, V2.String())
}
