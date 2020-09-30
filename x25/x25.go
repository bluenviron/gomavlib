// x25 implements the X25 hash.
package x25

// X25 is a X25 hasher.
type X25 struct {
	crc uint16
}

// New allocates a X25 hasher.
// X25 is the hash used to compute Frame checksums.
func New() *X25 {
	x := &X25{}
	x.Reset()
	return x
}

func (x *X25) Reset() {
	x.crc = 0xFFFF
}

func (x *X25) Size() int {
	return 2
}

func (x *X25) BlockSize() int {
	return 1
}

func (x *X25) Write(p []byte) (int, error) {
	for _, b := range p {
		tmp := uint16(b) ^ (x.crc & 0xFF)
		tmp ^= (tmp << 4)
		tmp &= 0xFF
		x.crc = (x.crc >> 8) ^ (tmp << 8) ^ (tmp << 3) ^ (tmp >> 4)
	}
	return len(p), nil
}

func (x *X25) Sum16() uint16 {
	return x.crc
}

func (x *X25) Sum(b []byte) []byte {
	return append(b, byte(x.crc), byte(x.crc>>8))
}
