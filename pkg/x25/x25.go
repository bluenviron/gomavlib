// Package x25 implements the X25 hash.
package x25

// X25 is the hash used to compute Frame checksums.
type X25 struct {
	crc uint16
}

// New allocates a X25.
func New() *X25 {
	x := &X25{}
	x.Reset()
	return x
}

// Reset resets the Hash to its initial state.
func (x *X25) Reset() {
	x.crc = 0xFFFF
}

// Size returns the number of bytes Sum will return.
func (x *X25) Size() int {
	return 2
}

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
func (x *X25) BlockSize() int {
	return 1
}

// Write adds more data to the running hash.
func (x *X25) Write(p []byte) (int, error) {
	for _, b := range p {
		tmp := uint16(b) ^ (x.crc & 0xFF)
		tmp ^= (tmp << 4)
		tmp &= 0xFF
		x.crc = (x.crc >> 8) ^ (tmp << 8) ^ (tmp << 3) ^ (tmp >> 4)
	}
	return len(p), nil
}

// Sum16 returns the curren thash.
func (x *X25) Sum16() uint16 {
	return x.crc
}

// Sum appends the current hash to b and returns the resulting slice.
func (x *X25) Sum(b []byte) []byte {
	return append(b, byte(x.crc), byte(x.crc>>8))
}
