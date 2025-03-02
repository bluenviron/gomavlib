# gomavlib

[![Test](https://github.com/bluenviron/gomavlib/workflows/test/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:test)
[![Lint](https://github.com/bluenviron/gomavlib/workflows/lint/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:lint)
[![Dialects](https://github.com/bluenviron/gomavlib/workflows/dialects/badge.svg)](https://github.com/bluenviron/gomavlib/actions?query=workflow:dialects)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluenviron/gomavlib)](https://goreportcard.com/report/github.com/bluenviron/gomavlib)
[![CodeCov](https://codecov.io/gh/bluenviron/gomavlib/branch/main/graph/badge.svg)](https://app.codecov.io/gh/bluenviron/gomavlib/tree/main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bluenviron/gomavlib/v3)](https://pkg.go.dev/github.com/bluenviron/gomavlib/v3#pkg-index)

gomavlib is a library that implements the Mavlink protocol (2.0 and 1.0) in the Go programming language. It can interact with Mavlink-capable devices through a serial port, UDP, TCP or a custom transport, and it can be used to power UGVs, UAVs, ground stations, monitoring systems or routers.

Mavlink is a lightweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by the most popular open-source flight controllers (Ardupilot and PX4).

This library powers the [**mavp2p**](https://github.com/bluenviron/mavp2p) router.

Features:

* Create Mavlink nodes able to communicate with other nodes.
  * Supported transports: serial, UDP (server, client or broadcast mode), TCP (server or client mode), custom reader/writer.
  * Emit heartbeats automatically.
  * Send automatic stream requests to Ardupilot devices (disabled by default).
  * Use both domain names and IPs.
* Decode and encode Mavlink v2.0 and v1.0.
  * Compute and validate checksums.
  * Support all v2 features: empty-byte truncation, signatures, message extensions.
* Use dialects in multiple ways.
  * Ready-to-use standard dialects are available in directory `dialects/`.
  * Custom dialects can be defined. Aa dialect generator is available in order to convert XML definitions into their Go representation.
  * Use no dialect at all. Messages can be routed without having their content decoded.

## Table of contents

* [Installation](#installation)
* [Examples](#examples)
* [API Documentation](#api-documentation)
* [Dialect generation](#dialect-generation)
* [Specifications](#specifications)
* [Links](#links)

## Installation

1. Install Go &ge; 1.21.

2. Create an empty folder, open a terminal in it and initialize the Go modules system:

   ```
   go mod init main
   ```

3. Download one of the [example files](#examples) and place it in the folder.

4. Compile and run:

   ```
   go run name-of-the-go-file.go
   ```

## Examples

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
* [dialect-absent](examples/dialect-absent/main.go)
* [dialect-custom](examples/dialect-custom/main.go)
* [events](examples/events/main.go)
* [router](examples/router/main.go)
* [router-edit](examples/router-edit/main.go)
* [serial-to-json](examples/serial-to-json/main.go)
* [stream-requests](examples/stream-requests/main.go)
* [read-writer](examples/read-writer/main.go)

## API Documentation

[Click to open the API Documentation](https://pkg.go.dev/github.com/bluenviron/gomavlib/v3#pkg-index)

## Dialect generation

Standard dialects are provided in the `pkg/dialects/` folder, but it's also possible to use custom dialects, that can be converted into Go files by running:

```
go install github.com/bluenviron/gomavlib/v3/cmd/dialect-import@latest
dialect-import my_dialect.xml
```

## Specifications

|name|area|
|----|----|
|[main website](https://mavlink.io/en/)|protocol|
|[packet format](https://mavlink.io/en/guide/serialization.html)|protocol|
|[common dialect](https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/common.xml)|dialects|
|[Golang project layout](https://github.com/golang-standards/project-layout)|project layout|

## Links

Related projects

* [mavp2p](https://github.com/bluenviron/mavp2p)

Other Go libraries

* [gobot](https://github.com/hybridgroup/gobot/tree/master/platforms/mavlink)
* [liamstask/go-mavlink](https://github.com/liamstask/go-mavlink)
* [ungerik/go-mavlink](https://github.com/ungerik/go-mavlink)
* [SpaceLeap/go-mavlink](https://github.com/SpaceLeap/go-mavlink)
* [MAVSDK-Go](https://github.com/mavlink/MAVSDK-Go)

Other non-Go libraries

* [official library (C)](https://github.com/mavlink/c_library_v2)
* [pymavlink (Python)](https://github.com/ArduPilot/pymavlink)
* [mavlink.net (C#)](https://github.com/asvol/mavlink.net)
* [rust-mavlink (Rust)](https://github.com/3drobotics/rust-mavlink)
* [node-mavlink (JS)](https://github.com/omcaree/node-mavlink)
