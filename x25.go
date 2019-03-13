package gomavlib

import (
	"hash"
)

// Hash16 is an interface modeled on the standard hash.Hash32 and hash.Hash64.
// it is implemented by 16-bits hash functions. In this library, it is used byX25.
type Hash16 interface {
	hash.Hash
	Sum16() uint16
}

type x25 struct {
	crc uint16
}

// NewX25 allocates a X25 hasher.
// X25 is the hash used to compute Frame checksums.
func NewX25() Hash16 {
	x := &x25{}
	x.Reset()
	return x
}

func (x *x25) Reset() {
	x.crc = 0xFFFF
}

func (x *x25) Size() int {
	return 2
}

func (x *x25) BlockSize() int {
	return 1
}

func (x *x25) Write(p []byte) (int, error) {
	for _, b := range p {
		tmp := uint16(b) ^ (x.crc & 0xFF)
		tmp ^= (tmp << 4)
		tmp &= 0xFF
		x.crc = (x.crc >> 8) ^ (tmp << 8) ^ (tmp << 3) ^ (tmp >> 4)
	}
	return len(p), nil
}

func (x *x25) Sum16() uint16 {
	return x.crc
}

func (x *x25) Sum(b []byte) []byte {
	return append(b, byte(x.crc), byte(x.crc>>8))
}
