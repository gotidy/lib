[![GoDev](https://img.shields.io/static/v1?label=godev&message=reference&color=00add8)][godev]

[godev]: https://pkg.go.dev/github.com/gotidy/lib/ptr

`ptr` contains functions for simplified creation of pointers from constants of basic types.

## Examples

### Getting pointer

This code:

```go
p := ptr.Of(10)
```

is the equivalent for:

```go
i := int(10)
p := &i  
```

### Getting value

```go

p := ptr.Of(10) 
v := ptr.Value(ptr.Of(10)) // This code returns 10.

p = nil
v = ptr.Value(p) // This code returns the default value 0.

v = ptr.ValueDef(p, 100) // This code returns default value 100.

```

## Documentation

[GoDoc](http://godoc.org/github.com/gotidy/ptr)

## License

[Apache 2.0](https://github.com/gotidy/lib/blob/master/LICENSE)
