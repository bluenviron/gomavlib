
# gomavlib

[![GoDoc](https://godoc.org/github.com/gswly/gomavlib?status.svg)](https://godoc.org/github.com/gswly/gomavlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/gswly/gomavlib)](https://goreportcard.com/report/github.com/gswly/gomavlib)

gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go programming language. It can power UGVs, UAVs, ground stations, monitoring systems or routers acting in a Mavlink network.

Mavlink is a lighweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by both of the most common open source drone softwares (Ardupilot and PX4).

## Features

* Decode and encode Mavlink v2.0 and v1.0. Supports checksums, empty-byte truncation (v2.0), signing (v2.0), message extensions (v2.0)
* Dialect is optional, the library can work with no dialects, with a user-defined dialect or with a standard dialect. Standard dialects are provided in directory `dialect/`, with no need for generation. A Dialect generator is provided anyway.
* Provides a ready-to-use Mavlink node with heartbeat emission and ability to communicate through multiple transports in parallel:
  * serial
  * UDP (server, client or broadcast mode)
  * TCP (server or client mode)
  * custom reader/writer
* UDP connections are tracked and removed when inactive
* Support both domain names and IPs
* Examples provided for every feature
* Comprehensive test suite

## Installation

Go &ge; 1.11 is required. If modules are enabled (i.e. there's a go.mod file in your project folder), it is enough to write the library name in the import section of the source files that are referring to it. Go will take care of downloading the needed files:
```go
import (
    ...
    "github.com/gswly/gomavlib"
)
```

If modules are not enabled, the library must be downloaded manually:
```
go get github.com/gswly/gomavlib
```

## Examples

* [serial_read_write](example/1serial_read_write.go)
* [udp_read_write](example/2udp_read_write.go)

## Documentation

https://godoc.org/github.com/gswly/gomavlib

## Testing

If you want to edit the library and test the results, unit tests can be run through:
```
make test
```

## Links

Protocol references
* https://mavlink.io/en/ (packet format: https://mavlink.io/en/guide/serialization.html)

Other Go libraries
* https://github.com/hybridgroup/gobot/tree/master/platforms/mavlink
* https://github.com/liamstask/go-mavlink
* https://github.com/ungerik/go-mavlink
* https://github.com/SpaceLeap/go-mavlink

Other non-Go libraries
* [C] https://github.com/mavlink/c_library_v2
* [Python] https://github.com/ArduPilot/pymavlink
* [Java] https://github.com/DrTon/jMAVlib
* [C#] https://github.com/asvol/mavlink.net
