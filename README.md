# FastID -- Snowflake ID generating in Go

[![Godoc](https://godoc.org/github.com/beinan/fastid?status.svg)](https://godoc.org/github.com/beinan/fastid)
[![Go Report Card](https://goreportcard.com/badge/github.com/beinan/fastid)](https://goreportcard.com/report/github.com/beinan/fastid)
[![Build Status](https://travis-ci.org/beinan/fastid.svg?branch=master)](https://travis-ci.org/beinan/fastid)
[![codecov](https://codecov.io/gh/beinan/fastid/branch/master/graph/badge.svg)](https://codecov.io/gh/beinan/fastid)

FastID is a pluggable unique ID generator in Go. 

* Under 64 bits (Long Integer)
* K-Ordered
* Lock-free (using atomic CAS)
* Decentralized

## Installation

```bash
go get github.com/beinan/fastid
```
## Quick Start
Generate an ID
```go

import (
  "fmt"
  "github.com/beinan/fastid"
)

func ExampleGenInt64ID() {
  id := fastid.CommonConfig.GenInt64ID()
  fmt.Printf("id generated: %v", id)
}
```

### Recommended Settings
 * 40 bits timestamp (34 years from 2018-06-01)
 * 16 bits machine ID (using lower 16 bits of IP v4 addr as default)
 * 7  bits sequence number
 
With this setting, FastID is able to generate 128(2^7) unique IDs per millisecond (1.048576 millisecond, 2^10 nanosecond).

### Customized Settings
See the examples in [GoDoc](https://godoc.org/github.com/beinan/fastid)

## Benchmarks
### Benchmark Settings
 * 40 bits timestamp (34 years from 2018-06-01)
 * 8  bits machine ID
 * 15 bits sequence number
 
```bash
go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/beinan/fastid
BenchmarkGenID-4        20000000                79.7 ns/op
BenchmarkGenIDP-4       20000000               141 ns/op
PASS
ok      github.com/beinan/fastid        4.779s
```
