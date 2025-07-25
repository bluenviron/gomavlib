//autogenerated:yes
//nolint:revive
package marsh

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
)

func TestDialect(t *testing.T) {
	d := &dialect.ReadWriter{Dialect: Dialect}
	err := d.Initialize()
	require.NoError(t, err)
}
