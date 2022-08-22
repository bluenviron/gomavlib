package main

import (
	"bytes"
	"fmt"

	"github.com/aler9/gomavlib/pkg/dialect"
	"github.com/aler9/gomavlib/pkg/dialects/ardupilotmega"
	"github.com/aler9/gomavlib/pkg/frame"
)

// if NewNode() is not flexible enough, the library provides a low-level
// frame reader and writer, that can be used with any kind of byte stream.

func main() {
	inBuf := bytes.NewBuffer(
		[]byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b"))
	outBuf := bytes.NewBuffer(nil)

	dialectRW, err := dialect.NewReadWriter(ardupilotmega.Dialect)
	if err != nil {
		panic(err)
	}

	rw, err := frame.NewReadWriter(
		frame.ReaderConf{
			Reader:    inBuf,
			DialectRW: dialectRW,
		},
		frame.WriterConf{
			Writer:      outBuf,
			DialectRW:   dialectRW,
			OutVersion:  frame.V2, // change to V1 if you're unable to communicate with the target
			OutSystemID: 10,
		})
	if err != nil {
		panic(err)
	}

	// read a message (wrapped in a frame)
	frame, err := rw.Read()
	if err != nil {
		panic(err)
	}

	fmt.Printf("decoded: %+v\n", frame)

	// write a message
	err = rw.WriteMessage(&ardupilotmega.MessageParamValue{
		ParamId:    "test_parameter",
		ParamValue: 123456,
		ParamType:  ardupilotmega.MAV_PARAM_TYPE_UINT32,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("encoded: %v\n", outBuf.Bytes())
}
