package uavionix

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/aler9/gomavlib/pkg/dialect"
)

func TestDialect(t *testing.T) {
	_, err := dialect.NewDecEncoder(Dialect)
	require.NoError(t, err)
}
