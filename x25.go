package gomavlib

// X25 is the hash used to compute Frame checksums.
// This structure implements the standard hash.Hash interface.
type X25 struct {
	crc uint16
}

// NewX25 allocates a X25.
func NewX25() *X25 {
	x := &X25{}
	x.Reset()
	return x
}

// Reset returns X25 to its initial state.
func (x *X25) Reset() {
	x.crc = 0xFFFF
}

// Size returns the size (in bytes) of the result of X25.
func (x *X25) Size() int {
	return 2
}

// BlockSize returns the preferred size of the input of X25.
func (x *X25) BlockSize() int {
	return 1
}

// Write insert some data for hash computation.
func (x *X25) Write(p []byte) (int, error) {
	for _, b := range p {
		tmp := uint16(b) ^ (x.crc & 0xFF)
		tmp ^= (tmp << 4)
		tmp &= 0xFF
		x.crc = (x.crc >> 8) ^ (tmp << 8) ^ (tmp << 3) ^ (tmp >> 4)
	}
	return len(p), nil
}

// Sum16 returns the hash of data written through Write.
func (x *X25) Sum16() uint16 {
	return x.crc
}

// Sum returns the hash of data written through Write.
func (x *X25) Sum(b []byte) []byte {
	return append(b, byte(x.crc), byte(x.crc>>8))
}
