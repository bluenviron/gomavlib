// Package frame contains frame definitions and a frame parser.
package frame

import (
	"bufio"

	"github.com/aler9/gomavlib/pkg/message"
)

// Frame is the interface implemented by frames of every supported version.
type Frame interface {
	// generates a clone of the frame.
	Clone() Frame

	// returns the system id of the author of the frame.
	GetSystemID() byte

	// returns the component id of the author of the frame.
	GetComponentID() byte

	// returns the message encapsuled in the frame.
	GetMessage() message.Message

	getChecksum() uint16
	decode(*bufio.Reader) error
	encodeTo([]byte, []byte) (int, error)
	genChecksum(byte) uint16
}
