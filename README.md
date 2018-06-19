# Fast ID

[![Godoc](https://godoc.org/github.com/beinan/fastid?status.svg)](https://godoc.org/github.com/beinan/fastid)
[![Build Status](https://travis-ci.org/beinan/fastid.svg?branch=master)](https://travis-ci.org/beinan/fastid)
[![codecov](https://codecov.io/gh/beinan/fastid/branch/master/graph/badge.svg)](https://codecov.io/gh/beinan/fastid)

## Benchmarks
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
