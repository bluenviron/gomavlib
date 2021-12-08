package udplistener

import (
	"bytes"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	testBuf1 := []byte("testing testing 1 2 3")
	testBuf2 := []byte("second part")

	l, err := New("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	defer l.Close()

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
}

func TestSamePacketMultipleReads(t *testing.T) {
	l, err := New("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		conn, err := l.Accept()
		require.NoError(t, err)
		defer conn.Close()

		buf := make([]byte, 256)

		for i := 0; i < 4; i++ {
			n, err := conn.Read(buf)
			require.NoError(t, err)
			require.Equal(t, 256, n)
		}
	}()

	conn, err := net.Dial("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	defer conn.Close()

	_, err = conn.Write(bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04}, 1024/4))
	require.NoError(t, err)

	wg.Wait()
}

func TestDeadline(t *testing.T) {
	l, err := New("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	defer l.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	var err1 error
	var err2 error

	go func() {
		defer wg.Done()

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
					err1 = err
					return
				}
				// second Read() must fail with Timeout
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					return
				}
				err1 = err
				return
			}
		}
	}()

	go func() {
		defer wg.Done()

		conn, err := net.Dial("udp4", "127.0.0.1:18456")
		require.NoError(t, err)
		defer conn.Close()

		_, err = conn.Write([]byte("a"))
		require.NoError(t, err)
	}()

	wg.Wait()
	require.NoError(t, err1)
	require.NoError(t, err2)
}

func TestDoubleClose(t *testing.T) {
	l, err := New("udp4", "127.0.0.1:18456")
	require.NoError(t, err)
	l.Close()
	l.Close()
}
