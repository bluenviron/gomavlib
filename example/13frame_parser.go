// +build ignore

package main

import (
	"fmt"
	"github.com/gswly/gomavlib"
	"github.com/gswly/gomavlib/dialects/ardupilotmega"
)

func main() {
	// if NewNode() is not enough, the library provides a low-level Mavlink
	// frame parser, that can be allocated with NewParser()
	parser, err := gomavlib.NewParser(gomavlib.ParserConf{
		Dialect: ardupilotmega.Dialect,
	})
	if err != nil {
		panic(err)
	}

	buf := []byte("\xfd\t\x01\x00\x00\x00\x00\x00\x00\x00\x04\x00\x00\x00\x01\x02\x03\x05\x03\xd9\xd1\x01\x02\x00\x00\x00\x00\x00\x0eG\x04\x0c\xef\x9b")

	// parse buf and obtain a frame
	frame, err := parser.Decode(buf, true, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("decoded: %+v\n", frame)
}
