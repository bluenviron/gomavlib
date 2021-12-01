package parser

import (
	"fmt"
)

// ReadError is the error returned in case of non-fatal parsing errors.
type ReadError struct {
	str string
}

func (e *ReadError) Error() string {
	return e.str
}

func newError(format string, args ...interface{}) *ReadError {
	return &ReadError{
		str: fmt.Sprintf(format, args...),
	}
}
