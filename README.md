
# gomavlib

[![Build Status](https://travis-ci.org/aler9/gomavlib.svg?branch=master)](https://travis-ci.org/aler9/gomavlib)
[![Go Report Card](https://goreportcard.com/badge/github.com/aler9/gomavlib)](https://goreportcard.com/report/github.com/aler9/gomavlib)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue)](https://pkg.go.dev/github.com/aler9/gomavlib?tab=doc)

gomavlib is a library that implements the Mavlink protocol (2.0 and 1.0) in the Go programming language. It can power UGVs, UAVs, ground stations, monitoring systems or routers, connected in a Mavlink network.

Mavlink is a lightweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by the most popular open-source flight controllers (Ardupilot and PX4).

This library powers the [**mavp2p**](https://github.com/aler9/mavp2p) router.

## Features

* Decodes and encodes Mavlink v2.0 and v1.0. Supports checksums, empty-byte truncation (v2.0), signatures (v2.0), message extensions (v2.0)
* Dialects are optional, the library can work with standard dialects (ready-to-use standard dialects are provided in directory `dialects/`), custom dialects or no dialects at all. In case of custom dialects, a dialect generator is available in order to convert XML definitions into their Go representation.
* Provides a high-level API (`Node`) with:
  * ability to communicate with multiple endpoints in parallel:
    * serial
    * UDP (server, client or broadcast mode)
    * TCP (server or client mode)
    * custom reader/writer
  * automatic heartbeat emission
  * automatic stream requests to Ardupilot devices (disabled by default)
* Provides a low-level API (`Parser`) with ability to decode/encode frames from/to a generic reader/writer
* UDP connections are tracked and removed when inactive
* Supports both domain names and IPs
* Examples provided for every feature, comprehensive test suite, continuous integration

## Installation

Go &ge; 1.12 is required, and modules must be enabled (there must be a `go.mod` file in your project folder, that can be created with the command `go mod init main`). To install the library, it is enough to write its name in the import section of the source files that will use it. Go will take care of downloading the needed files:
```go
import (
    "github.com/aler9/gomavlib"
)
```

## Examples

* [endpoint_serial](examples/endpoint_serial.go)
* [endpoint_udp_server](examples/endpoint_udp_server.go)
* [endpoint_udp_client](examples/endpoint_udp_client.go)
* [endpoint_udp_broadcast](examples/endpoint_udp_broadcast.go)
* [endpoint_tcp_server](examples/endpoint_tcp_server.go)
* [endpoint_tcp_client](examples/endpoint_tcp_client.go)
* [endpoint_custom](examples/endpoint_custom.go)
* [message_read](examples/message_read.go)
* [message_write](examples/message_write.go)
* [signature](examples/signature.go)
* [dialect_no](examples/dialect_no.go)
* [dialect_custom](examples/dialect_custom.go)
* [events](examples/events.go)
* [router](examples/router.go)
* [stream_requests](examples/stream_requests.go)
* [parser](examples/parser.go)

## Dialect generation

Standard dialects are provided in the `dialects/` folder, but it's also possible to use custom dialects, that can be converted into Go files by running:
```
go get github.com/aler9/gomavlib/commands/dialgen
dialgen my_dialect.xml > dialect.go
```

## Documentation

https://pkg.go.dev/github.com/aler9/gomavlib

## Testing

If you want to hack the library and test the results, unit tests can be launched with:
```
make test
```

## Links

Protocol documentation
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
