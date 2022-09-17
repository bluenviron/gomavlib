package gomavlib

import (
	"io"
	"testing"
)

type testLoopback chan []byte

func (ch testLoopback) Close() error {
	close(ch)
	return nil
}

func (ch testLoopback) Read(buf []byte) (int, error) {
	ret, ok := <-ch
	if !ok {
		return 0, errTerminated
	}
	n := copy(buf, ret)
	return n, nil
}

func (ch testLoopback) Write(buf []byte) (int, error) {
	ch <- buf
	return len(buf), nil
}

type testEndpoint struct {
	io.ReadCloser
	io.Writer
}

func TestEndpointCustom(t *testing.T) {
	l1 := make(testLoopback)
	l2 := make(testLoopback)
	doTest(t, EndpointCustom{&testEndpoint{l1, l2}},
		EndpointCustom{&testEndpoint{l2, l1}})
}
