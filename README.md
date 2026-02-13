# gomavlib

[![Test](https://github.com/bluenviron/gomavlib/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/bluenviron/gomavlib/actions/workflows/test.yml?query=branch%3Amain)
[![Lint](https://github.com/bluenviron/gomavlib/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/bluenviron/gomavlib/actions/workflows/lint.yml?query=branch%3Amain)
[![Dialects](https://github.com/bluenviron/gomavlib/actions/workflows/dialects.yml/badge.svg?branch=main)](https://github.com/bluenviron/gomavlib/actions/workflows/dialects.yml?query=branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluenviron/gomavlib)](https://goreportcard.com/report/github.com/bluenviron/gomavlib)
[![CodeCov](https://codecov.io/gh/bluenviron/gomavlib/branch/main/graph/badge.svg)](https://app.codecov.io/gh/bluenviron/gomavlib/tree/main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bluenviron/gomavlib/v3)](https://pkg.go.dev/github.com/bluenviron/gomavlib/v3#pkg-index)

gomavlib is a library that implements the Mavlink protocol (2.0 and 1.0) in the Go programming language. It can interact with Mavlink-capable devices through a serial port, UDP, TCP or a custom transport, and it can be used to power UGVs, UAVs, ground stations, monitoring systems or routers.

Mavlink is a lightweight and transport-independent protocol that is mostly used to communicate with unmanned ground vehicles (UGV) and unmanned aerial vehicles (UAV, drones, quadcopters, multirotors). It is supported by the most popular open-source flight controllers (Ardupilot and PX4).

This library powers the [**mavp2p**](https://github.com/bluenviron/mavp2p) router.

Features:

* Create Mavlink nodes able to communicate with other nodes.
  * Supported transports: serial, UDP (server, client or broadcast mode), TCP (server or client mode), custom reader/writer.
  * Support both domain names and IPs.
  * Emit heartbeats automatically.
  * Send automatic stream requests to Ardupilot devices (disabled by default).
* Decode and encode Mavlink v2.0 and v1.0.
  * Compute and validate checksums.
  * Support all v2 features: empty-byte truncation, signatures, message extensions.
* Use dialects in multiple ways.
  * Ready-to-use standard dialects are available in directory `dialects/`.
  * Custom dialects can be defined. Aa dialect generator is available in order to convert XML definitions into their Go representation.
  * Use no dialect at all. Messages can be routed without having their content decoded.
* Read and write telemetry logs (tlog)

## Table of contents

* [Installation](#installation)
* [Examples](#examples)
* [API Documentation](#api-documentation)
* [Dialect generation](#dialect-generation)
* [Specifications](#specifications)
* [Links](#links)

## Installation

1. Install Go &ge; 1.25.

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

* [node-endpoint-serial](examples/node-endpoint-serial/main.go)
* [node-endpoint-udp-server](examples/node-endpoint-udp-server/main.go)
* [node-endpoint-udp-client](examples/node-endpoint-udp-client/main.go)
* [node-endpoint-udp-broadcast](examples/node-endpoint-udp-broadcast/main.go)
* [node-endpoint-tcp-server](examples/node-endpoint-tcp-server/main.go)
* [node-endpoint-tcp-client](examples/node-endpoint-tcp-client/main.go)
* [node-endpoint-custom-client](examples/node-endpoint-custom-client/main.go)
* [node-endpoint-custom-server](examples/node-endpoint-custom-server/main.go)
* [node-message-read](examples/node-message-read/main.go)
* [node-message-write](examples/node-message-write/main.go)
* [node-command-microservice](examples/node-command-microservice/main.go)
* [node-signature](examples/node-signature/main.go)
* [node-dialect-absent](examples/node-dialect-absent/main.go)
* [node-dialect-custom](examples/node-dialect-custom/main.go)
* [node-events](examples/node-events/main.go)
* [node-router](examples/node-router/main.go)
* [node-router-edit](examples/node-router-edit/main.go)
* [node-serial-to-json](examples/node-serial-to-json/main.go)
* [node-stream-requests](examples/node-stream-requests/main.go)
* [frame-read-writer](examples/frame-read-writer/main.go)
* [telemetry-log](examples/telemetry-log/main.go)

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
