package gomavlib

// Version is a Mavlink version.
type Version int

const (
	// V1 is Mavlink 1.0
	V1 Version = 1

	// V2 is Mavlink 2.0
	V2 Version = 2
)

// String implements fmt.Stringer.
func (v Version) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}
