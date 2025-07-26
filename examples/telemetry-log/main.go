// Package main contains an example.
package main

import (
	"fmt"
	"os"

	"github.com/bluenviron/gomavlib/v3/pkg/dialect"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
	"github.com/bluenviron/gomavlib/v3/pkg/tlog"
)

// this example shows how to:
// 1) open a telemetry log file.
// 2) print every telemetry log entry present inside the file.

func main() {
	// open a telemetry log file.
	f, err := os.Open("my-telemetry-log.tlog")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// allocate dialect reader / writer.
	dialectRW := &dialect.ReadWriter{Dialect: ardupilotmega.Dialect}
	err = dialectRW.Initialize()
	if err != nil {
		panic(err)
	}

	// allocate telemetry log reader.
	dec := tlog.Reader{
		ByteReader: f,
		DialectRW:  dialectRW,
	}
	err = dec.Initialize()
	if err != nil {
		panic(err)
	}

	// print every telemetry log entry present inside the file.
	for {
		var entry *tlog.Entry
		entry, err = dec.Read()
		if err != nil {
			panic(err)
		}

		fmt.Printf("date: %s message: %+v\n", entry.Time, entry.Frame.GetMessage())
	}
}
