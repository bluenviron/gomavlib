package parser

// WriterOutVersion is a Mavlink version.
type WriterOutVersion int

const (
	// V1 is Mavlink 1.0
	V1 WriterOutVersion = 1

	// V2 is Mavlink 2.0
	V2 WriterOutVersion = 2
)

// String implements fmt.Stringer.
func (v WriterOutVersion) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}
