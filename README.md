# Semantic Versioning for Go

[![Build](https://img.shields.io/travis/hansrodtang/semver.svg?style=flat)](https://travis-ci.org/hansrodtang/semver) [![Coverage](https://img.shields.io/coveralls/hansrodtang/semver.svg?style=flat)](https://coveralls.io/r/hansrodtang/semver) [![Issues](https://img.shields.io/github/issues/hansrodtang/semver.svg?style=flat)](https://github.com/hansrodtang/semver/issues) [![Tip](https://img.shields.io/gratipay/hansrodtang.svg?style=flat)](https://gratipay.com/hansrodtang/)
[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](http://choosealicense.com/licenses/mit/)

A [Semantic Versioning](http://semver.org/) library for [Go](http://golang.org).

Covers version `2.0.0` of the semver specification.

Documentation on the syntax for the `Satifies()` method can be found  [here](https://www.npmjs.org/doc/misc/semver.html).


## Installation

```
  go get github.com/hansrodtang/semver
```
For those who prefer it, you can also use [gopkg.in](http://gopkg.in):

```
  go get gopkg.in/hansrodtang/semver.v0
```

## Usage

```go
import github.com/hansrodtang/semver

v1, error := semver.New("1.5.0")
// do something with error
if v1.Satisfies("^1.0.0") {
  // do something
}
```

## Benchmarks

Test | Iterations | Time
------------------------|-----------|------------
BenchmarkParseSimple    | 5000000   | 356 ns/op
BenchmarkParseComplex   | 1000000   | 2200 ns/op
BenchmarkCompareSimple  | 500000000 | 3.85 ns/op
BenchmarkCompareComplex	| 100000000	| 17.3 ns/op

Run the benchmarks yourself with:

```
go test github.com/hansrodtang/semver -bench=.
```

## Tests

Run the tests with:

```
go test -cover github.com/hansrodtang/semver
```

## License

This software is licensed under the [MIT license](LICENSE.md).
