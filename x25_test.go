package gomavlib

import (
	"github.com/stretchr/testify/require"
	"testing"
)

/* Test vectors generated with

( docker build - -t temp << EOF
FROM amd64/python:3-stretch
RUN apt update && apt install -y --no-install-recommends \
    git \
    gcc \
    python3-dev \
    python3-setuptools \
    python3-wheel \
    python3-pip \
    python3-future \
    python3-lxml \
    && pip3 install pymavlink
EOF
) && docker run --rm -it temp python3 -c \
"from pymavlink.mavutil import x25crc; print('%.4x' % x25crc('\x01').crc);"

*/

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
