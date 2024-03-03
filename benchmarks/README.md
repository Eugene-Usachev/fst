# Benchmarks
### Hash function is HS256

# Parse

<h2>speed in ns/op (less is better)</h2>

![Image Description](images/parse_speed_ns.jpg)

<h2>speed in total repetitions (more is better)</h2>

![Image Description](images/parse_speed_total.jpg)

<h2>memory in bytes/op (less is better)</h2>

![Image Description](images/parse_memory.jpg)

<h2>allocs in allocs/op (less is better)</h2>

![Image Description](images/parse_allocs.jpg)

# Generate

<h2>speed in ns/op (less is better)</h2>

![Image Description](images/gen_speed_ns.jpg)

<h2>speed in total repetitions (more is better)</h2>

![Image Description](images/gen_speed_total.jpg)

<h2>memory in bytes/op (less is better)</h2>

![Image Description](images/gen_memory.jpg)

<h2>allocs in allocs/op (less is better)</h2>

![Image Description](images/gen_allocs.jpg)

# Table

goos: linux

goarch: amd64

pkg: benches

cpu: Intel(R) Core(TM) i5-9400F CPU @ 2.90GHz

| Benchmark                            | Iterations | Time (ns/op) | Bytes (B/op) | Allocations (allocs/op) |
|--------------------------------------|------------|--------------|--------------|-------------------------|
| UintGen_GoJose-6                     | 158667     | 7477         | 6760         | 84                      |
| UintGen_GolangJWT-6                  | 365820     | 3130         | 1904         | 31                      |
| UintGen_JWT_GO-6                     | 392596     | 2947         | 1584         | 27                      |
| UintGen_JWT-6                        | 653710     | 1745         | 800          | 13                      |
| UintGen_EncodedFST-6                 | 2024858    | 586.5        | 176          | 4                       |
| UintGen_FST-6                        | 2462677    | 487.3        | 80           | 2                       |
| BigStringGen_GoJose-6                | 42236      | 26358        | 19944        | 90                      |
| BigStringGen_GolangJWT-6             | 122193     | 9862         | 9889         | 32                      |
| BigStringGen_JWT_GO-6                | 122626     | 9476         | 9568         | 27                      |
| BigStringGen_JWT-6                   | 156385     | 7554         | 6028         | 14                      |
| BigStringGen_EncodedFST-6            | 281076     | 4227         | 4001         | 4                       |
| BigStringGen_FST-6                   | 403830     | 2793         | 1184         | 2                       |
| UintParse_GoJose-6                   | 177871     | 6636         | 4512         | 66                      |
| UintParse_GolangJWT-6                | 296571     | 3878         | 2208         | 39                      |
| UintParse_JWT_GO-6                   | 276013     | 4162         | 2680         | 42                      |
| UintParse_JWT-6                      | 364107     | 3231         | 2336         | 29                      |
| UintParse_EncodedFST-6               | 2116108    | 563.2        | 80           | 2                       |
| UintParse_FST-6                      | 2589146    | 465.7        | 32           | 1                       |
| BigStringParse_GoJose-6              | 43207      | 27488        | 13664        | 68                      |
| BigStringParse_GolangJWT-6           | 90196      | 13119        | 6976         | 40                      |
| BigStringParse_JWT_GO-6              | 83995      | 14201        | 8968         | 42                      |
| BigStringParse_JWT-6                 | 74436      | 16061        | 7256         | 29                      |
| BigStringParse_EncodedFST-6          | 292818     | 3970         | 1184         | 2                       |
| BigStringParse_FST-6                 | 468292     | 2556         | 32           | 1                       |
| UintGen_GoJose_PARALLEL-6            | 274987     | 4282         | 6761         | 84                      |
| UintGen_GolangJWT_PARALLEL-6         | 768590     | 1537         | 1906         | 31                      |
| UintGen_JWT_GO_PARALLEL-6            | 881838     | 1233         | 1585         | 27                      |
| UintGen_JWT_PARALLEL-6               | 1896991    | 600.3        | 800          | 13                      |
| UintGen_EncodedFST_PARALLEL-6        | 6421455    | 172.1        | 176          | 4                       |
| UintGen_FST_PARALLEL-6               | 9339318    | 120.0        | 80           | 2                       |
| BigStringGen_GoJose_PARALLEL-6       | 109754     | 13013        | 19949        | 90                      |
| BigStringGen_GolangJWT_PARALLEL-6    | 237013     | 5208         | 9908         | 32                      |
| BigStringGen_JWT_GO_PARALLEL-6       | 244294     | 4668         | 9578         | 27                      |
| BigStringGen_JWT_PARALLEL-6          | 364998     | 3128         | 6035         | 14                      |
| BigStringGen_EncodedFST_PARALLEL-6   | 603237     | 1909         | 4004         | 4                       |
| BigStringGen_FST_PARALLEL-6          | 1233139    | 934.6        | 1186         | 2                       |
| UintParse_GoJose_PARALLEL-6          | 380606     | 3147         | 4512         | 66                      |
| UintParse_GolangJWT_PARALLEL-6       | 647278     | 1712         | 2208         | 39                      |
| UintParse_JWT_GO_PARALLEL-6          | 591739     | 1908         | 2680         | 42                      |
| UintParse_JWT_PARALLEL-6             | 792940     | 1467         | 2336         | 29                      |
| UintParse_EncodedFST_PARALLEL-6      | 8961590    | 137          | 80           | 2                       |
| UintParse_FST_PARALLEL-6             | 12636765   | 94           | 32           | 1                       |
| BigStringParse_GoJose_PARALLEL-6     | 110364     | 10881        | 13666        | 68                      |
| BigStringParse_GolangJWT_PARALLEL-6  | 252014     | 4798         | 6977         | 40                      |
| BigStringParse_JWT_GO_PARALLEL-6     | 203853     | 5771         | 8970         | 42                      |
| BigStringParse_JWT_PARALLEL-6        | 232568     | 5005         | 7256         | 29                      |
| BigStringParse_EncodedFST_PARALLEL-6 | 1117708    | 1069         | 1185         | 2                       |
| BigStringParse_FST_PARALLEL-6        | 2613806    | 465.9        | 32           | 1                       |