
# gomavlib

[![Test](https://github.com/aler9/gomavlib/workflows/test/badge.svg)](https://github.com/aler9/gomavlib/actions?query=workflow:test)
[![Lint](https://github.com/aler9/gomavlib/workflows/lint/badge.svg)](https://github.com/aler9/gomavlib/actions?query=workflow:lint)
[![Dialects](https://github.com/aler9/gomavlib/workflows/dialects/badge.svg)](https://github.com/aler9/gomavlib/actions?query=workflow:dialects)
[![Go Report Card](https://goreportcard.com/badge/github.com/aler9/gomavlib)](https://goreportcard.com/report/github.com/aler9/gomavlib)
[![CodeCov](https://codecov.io/gh/aler9/gomavlib/branch/main/graph/badge.svg)](https://codecov.io/gh/aler9/gomavlib/branch/main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/aler9/gomavlib)](https://pkg.go.dev/github.com/aler9/gomavlib#pkg-index)

gomavlib is a library that implements the Mavlink protocol (2.0 and 1.0) in the Go programming language. It can power UGVs, UAVs, ground stations, monitoring systems or routers, connected to other Mavlink-capable devices through a serial port, UDP, TCP or a custom transport.

Mavlink is a lightweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by the most popular open-source flight controllers (Ardupilot and PX4).

This library powers the [**mavp2p**](https://github.com/aler9/mavp2p) router.

Features:

* Decode and encode Mavlink v2.0 and v1.0. Supports checksums, empty-byte truncation (v2.0), signatures (v2.0), message extensions (v2.0).
* Dialects are optional, the library can work with standard dialects (ready-to-use standard dialects are provided in directory `dialects/`), custom dialects or no dialects at all. In case of custom dialects, a dialect generator is available in order to convert XML definitions into their Go representation.
* Create nodes able to communicate with multiple endpoints in parallel and with multiple transports:
  * serial
  * UDP (server, client or broadcast mode)
  * TCP (server or client mode)
  * custom reader/writer
* Emit heartbeats automatically
* Send automatic stream requests to Ardupilot devices (disabled by default)
* Support both domain names and IPs
* Examples provided for every feature, comprehensive test suite, continuous integration

## Table of contents

* [Installation](#installation)
* [API Documentation](#api-documentation)
* [Dialect generation](#dialect-generation)
* [Testing](#testing)
* [Links](#links)

## Installation

1. Install Go &ge; 1.16.

2. Create an empty folder, open a terminal in it and initialize the Go modules system:

   ```
   go mod init main
   ```

3. Download one of the example files and place it in the folder:

  * [endpoint-serial](examples/endpoint-serial/main.go)
  * [endpoint-udp-server](examples/endpoint-udp-server/main.go)
  * [endpoint-udp-client](examples/endpoint-udp-client/main.go)
  * [endpoint-udp-broadcast](examples/endpoint-udp-broadcast/main.go)
  * [endpoint-tcp-server](examples/endpoint-tcp-server/main.go)
  * [endpoint-tcp-client](examples/endpoint-tcp-client/main.go)
  * [endpoint-custom](examples/endpoint-custom/main.go)
  * [message-read](examples/message-read/main.go)
  * [message-write](examples/message-write/main.go)
  * [signature](examples/signature/main.go)
  * [dialect-no](examples/dialect-no/main.go)
  * [dialect-custom](examples/dialect-custom/main.go)
  * [events](examples/events/main.go)
  * [router](examples/router/main.go)
  * [stream-requests](examples/stream-requests/main.go)
  * [parser](examples/parser/main.go)

4. Compile and run

   ```
   go run name-of-the-go-file.go
   ```

## API Documentation

https://pkg.go.dev/github.com/aler9/gomavlib#pkg-index

## Dialect generation

Standard dialects are provided in the `pkg/dialects/` folder, but it's also possible to use custom dialects, that can be converted into Go files by running:

```
go get github.com/aler9/gomavlib/cmd/dialect-import
dialect-import my_dialect.xml
```

## Testing

If you want to hack the library and test the results, unit tests can be launched with:

```
make test
```

## Links

Related projects

* https://github.com/aler9/mavp2p

Protocol documentation

* main website https://mavlink.io/en/
* packet format https://mavlink.io/en/guide/serialization.html
* common dialect https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/common.xml

Other Go libraries

* https://github.com/hybridgroup/gobot/tree/master/platforms/mavlink
* https://github.com/liamstask/go-mavlink
* https://github.com/ungerik/go-mavlink
* https://github.com/SpaceLeap/go-mavlink
* https://github.com/mavlink/MAVSDK-Go

Other non-Go libraries

* [C] https://github.com/mavlink/c_library_v2
* [Python] https://github.com/ArduPilot/pymavlink
* [C#] https://github.com/asvol/mavlink.net
* [Rust] https://github.com/3drobotics/rust-mavlink
* [JS] https://github.com/omcaree/node-mavlink

Conventions

* https://github.com/golang-standards/project-layout
