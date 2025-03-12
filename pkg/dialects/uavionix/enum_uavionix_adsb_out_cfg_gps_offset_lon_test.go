//autogenerated:yes
//nolint:revive,govet,errcheck
package uavionix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnum_UAVIONIX_ADSB_OUT_CFG_GPS_OFFSET_LON(t *testing.T) {
	t.Run("zero", func(t *testing.T) {
		var e UAVIONIX_ADSB_OUT_CFG_GPS_OFFSET_LON
		e.UnmarshalText([]byte{})
		e.MarshalText()
		e.String()
	})

	t.Run("first entry", func(t *testing.T) {
		enc, err := UAVIONIX_ADSB_OUT_CFG_GPS_OFFSET_LON_NO_DATA.MarshalText()
		require.NoError(t, err)

		var dec UAVIONIX_ADSB_OUT_CFG_GPS_OFFSET_LON
		err = dec.UnmarshalText(enc)
		require.NoError(t, err)

		require.Equal(t, UAVIONIX_ADSB_OUT_CFG_GPS_OFFSET_LON_NO_DATA, dec)
	})
}
