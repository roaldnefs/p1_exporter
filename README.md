# P1 exporter

[![Travis CI](https://img.shields.io/travis/roaldnefs/p1_exporter.svg)](https://travis-ci.org/roaldnefs/p1_exporter)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/roaldnefs/p1_exporter)
[![Github All Releases](https://img.shields.io/github/downloads/roaldnefs/p1_exporter/total.svg)](https://github.com/roaldnefs/p1_exporter/releases)
[![GitHub](https://img.shields.io/github/license/roaldnefs/p1_exporter.svg)](https://github.com/roaldnefs/p1_exporter/blob/master/LICENSE)

Prometheus exporter for DSMR (Dutch Smart Meter Requirements) using the end-consumer (P1) interface.

* [Installation](README.md#installation)
     * [Binaries](README.md#binaries)
     * [Via Go](README.md#via-go)
* [Usage](README.md#usage)

## Installation

### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/roaldnefs/p1_exporter/releases).

### Via Go

```console
$ go get github.com/roaldnefs/p1_exporter
```

## Usage

```console
 $ p1_exporter--help
usage: p1_exporter --serial.port=SERIAL.PORT [<flags>]

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9602"  
                                 Address on which to expose metrics and web interface.
      --web.telemetry-path="/metrics"  
                                 Path under which to expose metrics.
      --serial.port=SERIAL.PORT  Serial port for the connection to the P1 interface.
```