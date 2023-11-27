# lib

Library based on generics. Slice, map helpers, set type, mathematic functions, pointer helpers.

[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev] [![Go Report Card](https://goreportcard.com/badge/github.com/gotidy/lib)][goreport]

<!-- [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) -->

[godev]: https://pkg.go.dev/github.com/gotidy/lib
[goreport]: https://goreportcard.com/report/github.com/gotidy/lib

## Installation

Required 1.18 or later version of Go.

```sh
go get -u github.com/gotidy/lib
```

## Documentation

### [Collections](collections/README.md)

- [Map](collections/maps/README.md) Contains map helpers.
- [Slice](collections/slice/README.md) Contains slice helpers.
- [Set](collections/set/README.md) Realize `Set` type.

### [Ptr](ptr/README.md)

Contains functions for simplified creation of pointers from constants of basic types.

### [Conversions](conversions/README.md)

Contains conversion helpers.

### [Conditions](conditions/README.md)

Contains functions that simplify inline conditions.

### [OneOf](oneof/README.md)

Realize `OneOf` type.

### [Math](math/README.md)

Contains mathematic functions.

### [Constraints](constraints/README.md)

Contains constraints types.

## License

[Apache 2.0](https://github.com/gotidy/lib/blob/master/LICENSE)
