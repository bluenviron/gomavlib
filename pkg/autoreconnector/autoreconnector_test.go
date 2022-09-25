package autoreconnector

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReconnect(t *testing.T) {
	ln, err := net.Listen("tcp", "localhost:6657")
	require.NoError(t, err)
	defer ln.Close()

	go func() {
		for i := 0; i < 2; i++ {
			conn, err := ln.Accept()
			require.NoError(t, err)

			_, err = conn.Write([]byte{0x05 + byte(i)})
			require.NoError(t, err)

			conn.Close()
		}
	}()

	a := New(
		func(ctx context.Context) (io.ReadWriteCloser, error) {
			return (&net.Dialer{}).DialContext(ctx, "tcp", "localhost:6657")
		},
	)
	defer a.Close()

	p := make([]byte, 1)
	n, err := a.Read(p)
	require.NoError(t, err)
	require.Equal(t, []byte{0x05}, p[:n])

	p = make([]byte, 1)
	n, err = a.Read(p)
	require.NoError(t, err)
	require.Equal(t, []byte{0x06}, p[:n])
}

func TestCloseWhileWorking(t *testing.T) {
	for _, ca := range []string{"read", "write"} {
		t.Run(ca, func(t *testing.T) {
			ln, err := net.Listen("tcp", "localhost:6657")
			require.NoError(t, err)
			defer ln.Close()

			serverDone := make(chan struct{})

			go func() {
				conn, err := ln.Accept()
				require.NoError(t, err)

				b := make([]byte, 1)
				if ca == "read" {
					conn.Read(b)
				} else {
					conn.Write(b)
				}

				conn.Close()
				close(serverDone)
			}()

			a := New(
				func(ctx context.Context) (io.ReadWriteCloser, error) {
					return (&net.Dialer{}).DialContext(ctx, "tcp", "localhost:6657")
				},
			)

			workDone := make(chan struct{})

			go func() {
				defer close(workDone)

				p := make([]byte, 1)
				if ca == "read" {
					a.Read(p)
				} else {
					a.Write(p)
				}
			}()

			time.Sleep(500 * time.Millisecond)
			a.Close()
			<-workDone
			<-serverDone
		})
	}
}
