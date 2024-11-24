// Package frame contains frame definitions and a frame parser.
package frame

import (
	"bufio"

	"github.com/chrisdalke/gomavlib/v3/pkg/message"
)

// Frame is the interface implemented by frames of every supported version.
type Frame interface {
	// returns the system id of the author of the frame.
	GetSystemID() byte

	// returns the component id of the author of the frame.
	GetComponentID() byte

	// returns the sequence number in the frame
	GetSequenceNumber() byte

	// returns the message wrapped in the frame.
	GetMessage() message.Message

	// returns the checksum of the frame.
	GetChecksum() uint16

	// generates the checksum of the frame.
	GenerateChecksum(byte) uint16

	decode(*bufio.Reader) error
	encodeTo([]byte, []byte) (int, error)
}
