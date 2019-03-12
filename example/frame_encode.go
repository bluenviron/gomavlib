// +build ignore

package main

import (
	"fmt"
)

func main() {
	frame := &gomavlib.FrameV2{
		SequenceId:  0x27,
		SystemId:    0x01,
		ComponentId: 0x02,
		Message: &ardupilotmega.MessageChangeOperatorControl{
			TargetSystem:   1,
			ControlRequest: 1,
			Version:        1,
			Passkey:        "testing",
		},
		Checksum: 0x66e5,
	}

}
