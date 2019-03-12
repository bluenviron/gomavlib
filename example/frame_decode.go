// +build ignore

package main

import (
	"fmt"
)

func main() {
	[]byte("\xfd\t\x00\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\x3e\x29")

	fmt.Printf("Frame version: %v\n", frame.GetVersion())
	fmt.Printf("Frame checksum: %v\n", frame.GetChecksum())
	fmt.Printf("Message type: %T\n", frame.GetMessage())
	fmt.Printf("Message content: %+v\n", frame.GetMessage())
}
