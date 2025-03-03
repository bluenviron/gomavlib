package main

import (
	"bytes"
	"io"
	"log"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/frame"
	"github.com/bluenviron/gomavlib/v3/pkg/streamwriter"
)

// When Node is not flexible enough, the library provides a low-level
// frame reader and writer, that can be used with any kind of byte stream.

// this example shows how to:
// 1) allocate the low-level frame.ReadWriter around a io.ReadWriter.
// 2) read a frame, that contains a message.
// 3) write a message, that is automatically wrapped in a frame.

type readWriter struct {
	io.Reader
	io.Writer
}

func main() {
	inBuf := bytes.NewBuffer(
		[]byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b"))
	outBuf := bytes.NewBuffer(nil)

	// allocate dialect reader / writer.
	dialectRW := &dialect.ReadWriter{Dialect: ardupilotmega.Dialect}
	err := dialectRW.Initialize()
	if err != nil {
		panic(err)
	}

	// allocate frame.ReadWriter around a io.ReadWriter.
	frameRW := &frame.ReadWriter{
		ByteReadWriter: &readWriter{
			Reader: inBuf,
			Writer: outBuf,
		},
		DialectRW: dialectRW,
	}
	err = frameRW.Initialize()
	if err != nil {
		panic(err)
	}

	// allocate streamwriter.Writer around frame.ReadWriter.
	streamWriter := &streamwriter.Writer{
		FrameWriter: frameRW.Writer,
		Version:     streamwriter.V2, // change to V1 if you're unable to communicate with the target
		SystemID:    10,
	}
	err = streamWriter.Initialize()
	if err != nil {
		panic(err)
	}

	// read a frame, that contains a message
	frame, err := frameRW.Read()
	if err != nil {
		panic(err)
	}

	log.Printf("decoded: %+v\n", frame)

	// write a message, that is automatically wrapped in a frame
	err = streamWriter.Write(&ardupilotmega.MessageParamValue{
		ParamId:    "test_parameter",
		ParamValue: 123456,
		ParamType:  ardupilotmega.MAV_PARAM_TYPE_UINT32,
	})
	if err != nil {
		panic(err)
	}

	log.Printf("encoded: %v\n", outBuf.Bytes())
}
