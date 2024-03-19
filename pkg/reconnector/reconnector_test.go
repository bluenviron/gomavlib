package reconnector

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type dummyRWC struct {
	bytes.Buffer
	closed bool
}

func (dummyRWC) Close() error {
	return nil
}

func (d *dummyRWC) Read(p []byte) (int, error) {
	return d.Buffer.Read(p)
}

func (d *dummyRWC) Write(p []byte) (int, error) {
	if d.closed {
		return 0, fmt.Errorf("closed")
	}
	return d.Buffer.Write(p)
}

func TestReconnector(t *testing.T) {
	var buf dummyRWC

	r := New(
		func(_ context.Context) (io.ReadWriteCloser, error) {
			return &buf, nil
		},
	)

	conn, ok := r.Reconnect()
	require.Equal(t, true, ok)

	buf.Buffer.Write([]byte{1})

	recv := make([]byte, 1)
	_, err := conn.Read(recv)
	require.NoError(t, err)
	require.Equal(t, byte(1), recv[0])

	_, err = conn.Read(recv)
	require.Equal(t, io.EOF, err)

	buf.closed = true
	_, err = conn.Write(recv)
	require.Error(t, err)

	_, ok = r.Reconnect()
	require.Equal(t, true, ok)

	r.Close()

	_, ok = r.Reconnect()
	require.Equal(t, false, ok)
}
