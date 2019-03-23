
# gomavlib

[![GoDoc](https://godoc.org/github.com/gswly/gomavlib?status.svg)](https://godoc.org/github.com/gswly/gomavlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/gswly/gomavlib)](https://goreportcard.com/report/github.com/gswly/gomavlib)
[![Build Status](https://travis-ci.org/gswly/gomavlib.svg?branch=master)](https://travis-ci.org/gswly/gomavlib)

gomavlib is a library that implements Mavlink 2.0 and 1.0 in the Go programming language. It can power UGVs, UAVs, ground stations, monitoring systems or routers acting in a Mavlink network.

Mavlink is a lighweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by both of the most common open-source flight controllers (Ardupilot and PX4).

This library powers the [**mavp2p**](https://github.com/gswly/mavp2p) router.

## Features

* Decode and encode Mavlink v2.0 and v1.0. Supports checksums, empty-byte truncation (v2.0), signatures (v2.0), message extensions (v2.0)
* Dialect is optional, the library can work a standard dialect, a custom dialect or no dialect at all. Standard dialects are provided in directory `dialects/`, with no need for generation. A Dialect generator is provided anyway.
* Provides a high-level node with ability to communicate with multiple endpoints in parallel:
  * serial
  * UDP (server, client or broadcast mode)
  * TCP (server or client mode)
  * custom reader/writer
* Provides a low-level parser with ability to decode/encode frames from/to a reader/writer
* UDP connections are tracked and removed when inactive
* Support both domain names and IPs
* Examples provided for every feature
* Comprehensive test suite

## Installation

Go &ge; 1.11 is required. If modules are enabled (i.e. there's a go.mod file in your project folder), it is enough to write the library name in the import section of the source files that are referring to it. Go will take care of downloading the needed files:
```go
import (
    "github.com/gswly/gomavlib"
)
```

If modules are not enabled, the library must be downloaded manually:
```
go get github.com/gswly/gomavlib
```

## Examples

* [endpoint_serial](example/1endpoint_serial.go)
* [endpoint_udp_server](example/2endpoint_udp_server.go)
* [endpoint_udp_client](example/3endpoint_udp_client.go)
* [endpoint_udp_broadcast](example/4endpoint_udp_broadcast.go)
* [endpoint_tcp_server](example/5endpoint_tcp_server.go)
* [endpoint_tcp_client](example/6endpoint_tcp_client.go)
* [endpoint_custom](example/7endpoint_custom.go)
* [message_write](example/8message_write.go)
* [message_signature](example/9message_signature.go)
* [dialect_no](example/10dialect_no.go)
* [dialect_custom](example/11dialect_custom.go)
* [events](example/12events.go)
* [router](example/13router.go)
* [frame_parser](example/14frame_parser.go)

## Documentation

https://godoc.org/github.com/gswly/gomavlib

## Dialect generation

Although standard dialects are provided in the `dialects/` folder, dialect definitions in XML format can be converted into Go files by running the `dialgen` utility:
```
go get github.com/gswly/gomavlib/dialgen
dialgen [path_or_url_to_xml_definition]
```

## Testing

If you want to edit the library and test the results, unit tests can be run through:
```
make test
```

## Links

Protocol references
* main website https://mavlink.io/en/
* packet format https://mavlink.io/en/guide/serialization.html
* common dialect https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/common.xml

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
* [Rust] https://github.com/3drobotics/rust-mavlink
* [JS] https://github.com/omcaree/node-mavlink
