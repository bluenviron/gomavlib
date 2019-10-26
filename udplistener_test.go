package gomavlib

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUdpListener(t *testing.T) {
	testBuf1 := []byte("testing testing 1 2 3")
	testBuf2 := []byte("second part")

	l, err := newUdpListener("udp4", "127.0.0.1:18456")
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()

		for i := 0; i < 2; i++ {
			conn, err := l.Accept()
			require.NoError(t, err)

			go func() {
				defer wg.Done()
				defer conn.Close()

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				require.NoError(t, err)
				require.Equal(t, len(testBuf1), n)
				require.Equal(t, testBuf1, buf[:n])

				n, err = conn.Write(testBuf2)
				require.NoError(t, err)
				require.Equal(t, len(testBuf2), n)
			}()
		}
	}()

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial("udp4", "127.0.0.1:18456")
			require.NoError(t, err)
			defer conn.Close()

			n, err := conn.Write(testBuf1)
			require.NoError(t, err)
			require.Equal(t, len(testBuf1), n)

			buf := make([]byte, 1024)
			n, err = conn.Read(buf)
			require.NoError(t, err)
			require.Equal(t, len(testBuf2), n)
			require.Equal(t, testBuf2, buf[:n])
		}()
	}

	wg.Wait()
	l.Close()
}

func TestUdpListenerDeadline(t *testing.T) {
	l, err := newUdpListener("udp4", "127.0.0.1:18456")
	require.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer l.Close()

		conn, err := l.Accept()
		require.NoError(t, err)
		defer conn.Close()

		for i := 0; i < 2; i++ {
			err = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			require.NoError(t, err)

			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err != nil {
				// accept first Read()
				if i == 0 {
					t.Fatal(err)
				}
				// second Read() must fail with Timeout
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					return
				}
				t.Fatal(err)
			}
		}
	}()

	go func() {
		defer wg.Done()

		conn, err := net.Dial("udp4", "127.0.0.1:18456")
		require.NoError(t, err)
		defer conn.Close()

		conn.Write([]byte("a"))
	}()

	wg.Wait()
}

func TestUdpListenerDoubleClose(t *testing.T) {
	l, err := newUdpListener("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	l.Close()
	l.Close()
}
