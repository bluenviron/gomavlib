package gomavlib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestX25(t *testing.T) {
	ins := [][]byte{
		{},
		[]byte("\x01"),
		[]byte("hello world"),
		[]byte("The quick brown fox jumps over the lazy dog"),
	}
	outs := []uint16{
		0xFFFF,
		0x1e0e,
		0x51f9,
		0x6ca7,
	}

	for i, in := range ins {
		h := NewX25()
		h.Write(in)
		out := h.Sum16()
		require.Equal(t, outs[i], out)
	}
}
