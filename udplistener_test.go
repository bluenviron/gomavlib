package gomavlib

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestUdpListener(t *testing.T) {
	testBuf1 := []byte("testing testing 1 2 3")
	testBuf2 := []byte("second part")

	l, err := newUdpListener("udp4", ":18456")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()

		for i := 0; i < 2; i++ {
			conn, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}

			go func() {
				defer wg.Done()
				defer conn.Close()

				err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				if err != nil {
					t.Fatal(err)
				}

				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				if err != nil {
					t.Fatal(err)
				}
				if n != len(testBuf1) {
					t.Fatal("length different")
				}
				if reflect.DeepEqual(buf[:n], testBuf1) == false {
					t.Fatalf("buffer different")
				}

				n, err = conn.Write(testBuf2)
				if err != nil {
					t.Fatal(err)
				}
				if n != len(testBuf2) {
					t.Fatalf("unexpected length")
				}
			}()
		}
	}()

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()

			conn, err := net.Dial("udp4", "127.0.0.1:18456")
			if err != nil {
				t.Fatal(err)
			}

			n, err := conn.Write(testBuf1)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(testBuf1) {
				t.Fatalf("unexpected length")
			}

			err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			if err != nil {
				t.Fatal(err)
			}

			buf := make([]byte, 1024)
			n, err = conn.Read(buf)
			if err != nil {
				t.Fatal(err)
			}
			if n != len(testBuf2) {
				t.Fatal("length different")
			}
			if reflect.DeepEqual(buf[:n], testBuf2) == false {
				t.Fatalf("buffer different")
			}
		}()
	}

	wg.Wait()
	l.Close()
}

func TestUdpListenerDeadline(t *testing.T) {
	l, err := newUdpListener("udp4", ":18456")
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer l.Close()

		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i < 2; i++ {
			err = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
			if err != nil {
				t.Fatal(err)
			}

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
		if err != nil {
			t.Fatal(err)
		}
		defer conn.Close()

		conn.Write([]byte("a"))
	}()

	wg.Wait()
}
