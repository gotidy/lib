[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]

[godev]: https://pkg.go.dev/github.com/gotidy/lib/ptr

`ptr` contains functions for simplified creation of pointers from constants of basic types.

## Examples

This code:

```go
p := ptr.Of(10)
```

is the equivalent for:

```go
i := int(10)
p := &i  
```

## Documentation

[GoDoc](http://godoc.org/github.com/gotidy/ptr)

## License

[Apache 2.0](https://github.com/gotidy/lib/blob/master/LICENSE)
