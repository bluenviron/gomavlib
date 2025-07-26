package frame

import (
	"bytes"
	"testing"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
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

func TestReadWriterNewErrors(t *testing.T) {
	_, err := NewReadWriter(ReadWriterConf{
		OutVersion:  V2,
		OutSystemID: 1,
	})
	require.EqualError(t, err, "BufByteReader not provided")
}

func FuzzReadWriter(f *testing.F) {
	for _, ca := range casesReadWrite {
		f.Add(ca.raw, false, false)
	}

	dialectRW := &dialect.ReadWriter{Dialect: ardupilotmega.Dialect}
	err := dialectRW.Initialize()
	if err != nil {
		panic(err)
	}

	f.Fuzz(func(t *testing.T, a []byte, k bool, v2 bool) {
		var key *V2Key
		if k {
			key = NewV2Key(bytes.Repeat([]byte("\x4F"), 32))
		}

		var outv WriterOutVersion
		if v2 {
			outv = V2
		} else {
			outv = V1
		}

		buf := bytes.NewBuffer(a)
		rw := &ReadWriter{
			ByteReadWriter: buf,
			DialectRW:      dialectRW,
			InKey:          key,
			OutVersion:     outv,
			OutSystemID:    1,
		}
		err2 := rw.Initialize()
		require.NoError(t, err2)

		fr, err2 := rw.Read()
		if err2 != nil {
			return
		}

		err2 = rw.Write(fr)
		require.NoError(t, err2)
	})
}
