package frame

// WriterOutVersion is a Mavlink version.
//
// Deprecated: replaced by streamwriter.Version.
type WriterOutVersion int

const (
	// V1 is Mavlink 1.0
	//
	// Deprecated: switch to streamwriter.
	V1 WriterOutVersion = 1

	// V2 is Mavlink 2.0
	//
	// Deprecated: switch to streamwriter.
	V2 WriterOutVersion = 2
)

// String implements fmt.Stringer.
func (v WriterOutVersion) String() string {
	if v == V1 {
		return "V1"
	}
	return "V2"
}
