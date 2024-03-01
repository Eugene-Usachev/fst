# fst

fst is a high-performance, lightweight library for generating and parsing Fast Signed Token (FST). 
FST provides an alternative to JSON-based tokens and allows you to store any information that can be 
represented as `[]byte`. You can use FST for the same purposes as JWT.

## Description

fst is designed to be efficient and lightweight, making it ideal for applications that require fast token generation
and parsing. With its optimized algorithms and data structures, fst minimizes memory usage and maximizes performance.

## Performance

fst excels in terms of performance, especially when compared to traditional token formats like JSON Web Tokens (JWT).
By leveraging its unique token structure and optimized parsing algorithms, fst significantly reduces the overhead
associated with token generation and parsing.

To demonstrate the performance benefits of fst, we conducted a series of tests using various tokens sizes.
The test results and source code for the tests are available on https://github.com/Eugene-Usachev/fst/tree/main/benchmarks.

## Installation

Install fst with the go get command:

`go get github.com/Eugene-Usachev/fst`

## Example

You can see the examples in the `example` folder.

First, you need to create a `Converter'. You can do it like this

```go
converter := fst.NewConverter(&fst.ConverterConfig{
    SecretKey:   []byte(`secret`),
    Postfix:     []byte(`postfix`),
    HashType:    sha256.New,
})
```

Then you can create a token using the `NewToken`

```go
token := converter.NewToken([]byte(`token`))
```

To parse tokens, you can use the `ParseToken`

```go
value, err := converter.ParseToken(token)
```

If you want to set expiration time, create a new converter

```go
converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
    SecretKey:          []byte(`secret`),
    Postfix:            nil,
    ExpirationTime:     time.Minute * 5,
    HashType:           sha256.New,
})
```

### Attention, please!
For work with the browser and HTTP in general, use EncodedConverter!
But if you don't need it, you can use Converter instead, as it is faster and more lightweight.

To create EncodedConverter call NewEncodedConverter with ConverterConfig
```go
encodedConverter := converter := fst.NewEncodedConverter(&fst.ConverterConfig{
    SecretKey: []byte(`secret`),
    HashType:  sha256.New,
})

converterWithExpirationTime := fst.NewEncodedConverter(&fst.ConverterConfig{
    SecretKey:      []byte(`secret`),
    ExpirationTime: time.Minute * 5,
    HashType:       sha256.New,
})
```

## License

The `fst` library is released under the MIT License.
