# fst

fst is a high-performance, low-memory library for generating and parsing Fast Signed Token (FST). 
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
The results are as follows:

<h2>Parse</h2>

<h3>speed in ns/op (less is better)</h3>

![Image Description](benchmarks/images/parse_speed_ns.jpg)

<h3>speed in total repetitions (more is better)</h3>

![Image Description](benchmarks/images/parse_speed_total.jpg)

<h3>memory in bytes/op (less is better)</h3>

![Image Description](benchmarks/images/parse_memory.jpg)

<h3>allocs in allocs/op (less is better)</h3>

![Image Description](benchmarks/images/parse_allocs.jpg)

<h2>Generate</h2>

<h3>speed in ns/op (less is better)</h3>

![Image Description](benchmarks/images/gen_speed_ns.jpg)

<h3>speed in total repetitions (more is better)</h3>

![Image Description](benchmarks/images/gen_speed_total.jpg)

<h3>memory in bytes/op (less is better)</h3>

![Image Description](benchmarks/images/gen_memory.jpg)

<h3>allocs in allocs/op (less is better)</h3>

![Image Description](benchmarks/images/gen_allocs.jpg)

<h2>Parallel parse</h2>

<h3>speed in ns/op (less is better)</h3>

![Image Description](benchmarks/images/parallel_parse_ns.jpg)

<h3>speed in total repetitions (more is better)</h3>

![Image Description](benchmarks/images/parallel_parse_total.jpg)

<h2>Parallel generate</h2>

<h3>speed in ns/op (less is better)</h3>

![Image Description](benchmarks/images/parallel_gen_ns.jpg)

<h3>speed in total repetitions (more is better)</h3>

![Image Description](benchmarks/images/parallel_gen_total.jpg)


To learn more about benchmarks, you can visit the `benchmarks` folder. There you will find the source code and data in the form of a table.

## Installation

Install fst with the go get command:

`go get github.com/Eugene-Usachev/fst`

## Example

You can see the examples in the `example` folder.

First you need to create a `Converter'. You can do it like this

```go
converter := fst.NewConverter(&fst.ConverterConfig{
    SecretKey:   []byte(`secret`),
    Postfix:     nil,
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

If you want to set expiration time create new converter

```go
converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
    SecretKey:          []byte(`secret`),
    Postfix:            nil,
    ExpirationTime:     time.Minute * 5,
    HashType:           sha256.New,
    WithExpirationTime: true,
})
```
## License

The `fst` library is released under the MIT License.
