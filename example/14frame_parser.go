// +build ignore

package main

import (
	"bytes"
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	inBuf := bytes.NewBuffer(
		[]byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b"))
	outBuf := bytes.NewBuffer(nil)

	// if NewNode() is not flexible enough, the library provides a low-level Mavlink
	// frame parser, that can be allocated with NewParser().
	parser, err := gomavlib.NewParser(gomavlib.ParserConf{
		Reader:      inBuf,
		Writer:      outBuf,
		Dialect:     ardupilotmega.Dialect,
		OutSystemId: 10,
	})
	if err != nil {
		panic(err)
	}

	// parse buffer and obtain a frame
	frame, err := parser.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("decoded: %+v\n", frame)

	// encode a frame
	frame = &gomavlib.FrameV2{
		Message: &ardupilotmega.MessageParamValue{
			ParamId:    "test_parameter",
			ParamValue: 123456,
			ParamType:  uint8(ardupilotmega.MAV_PARAM_TYPE_UINT32),
		},
	}
	err = parser.Write(frame, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("encoded: %v\n", outBuf.Bytes())
}
