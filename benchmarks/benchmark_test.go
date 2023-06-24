package benchmarks

import (
	v5 "benches/goJose"
	v4 "benches/golangJwt"
	v3 "benches/jwt"
	v2 "benches/jwt-go"
	"github.com/Eugene-Usachev/fst"
	"strconv"
	"strings"
	"testing"
)

var (
	key1 = "key1"
)

func U2B(u uint) []byte {
	return []byte(strconv.FormatUint(uint64(u), 10))
}

var (
	bkey1 = []byte(key1)
	id    = uint(1)

	fstConverterA = fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	fstConverterR = fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})

	message1 = func() string {
		b := strings.Builder{}
		for i := 0; i < 1000; i++ {
			b.WriteString("a")
		}
		str := b.String()
		return str
	}()
	uintTokenV2, _      = v2.NewAccessToken(id, bkey1)
	uintTokenV3, _      = v3.NewAccessToken(id)
	uintTokenV4, _      = v4.NewAccessToken(id, bkey1)
	uintTokenV5, _      = v5.NewAccessToken(id)
	uintTokenFST        = fstConverterA.NewToken(U2B(id))
	bigStringTokenV2, _ = v2.NewRefreshToken(bkey1, message1)
	bigStringTokenV3, _ = v3.NewRefreshToken(message1)
	bigStringTokenV4, _ = v4.NewRefreshToken(message1, bkey1)
	bigStringTokenV5, _ = v5.NewRefreshToken(message1)
	bigStringTokenFST   = fstConverterR.NewToken([]byte(message1))
)

func BenchmarkUintGen_GoJose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v5.NewAccessToken(id)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkUintGen_GolangJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v4.NewAccessToken(id, bkey1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkUintGen_JWT_GO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v2.NewAccessToken(id, bkey1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkUintGen_JWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v3.NewAccessToken(id)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkUintGen_FST(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	bid := U2B(id)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token := converter.NewToken(bid)
		if len(token) < 1 {
		}
	}
}

func BenchmarkBigStringGen_GoJose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v5.NewRefreshToken(message1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkBigStringGen_GolangJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v4.NewRefreshToken(message1, bkey1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkBigStringGen_JWT_GO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v2.NewRefreshToken(bkey1, message1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkBigStringGen_JWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		token, err := v3.NewRefreshToken(message1)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkBigStringGen_FST(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	bmessage := []byte(message1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token := converter.NewToken(bmessage)
		if len(token) < 1 {
		}
	}
}

func BenchmarkUintParse_GoJose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idT, err := v5.ParseAccessToken(uintTokenV5)
		if err != nil && idT < 1 {
		}
	}
}

func BenchmarkUintParse_GolangJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idT, err := v4.ParseAccessToken(uintTokenV4, bkey1)
		if err != nil && idT < 1 {
		}
	}
}

func BenchmarkUintParse_JWT_GO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idT, err := v2.ParseAccessToken(uintTokenV2, bkey1)
		if err != nil && idT < 1 {
		}
	}
}

func BenchmarkUintParse_JWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		idT, err := v3.ParseAccessToken(uintTokenV3)
		if err != nil && idT < 1 {
		}
	}
}

func BenchmarkUintParse_FST(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token, err := converter.ParseToken(uintTokenFST)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkBigStringParse_GoJose(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pass, err := v5.ParseRefreshToken(bigStringTokenV5)
		if err != nil && len(pass) < 1 {
		}
	}
}

func BenchmarkBigStringParse_GolangJWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pass, err := v4.ParseRefreshToken(bigStringTokenV4, bkey1)
		if err != nil && len(pass) < 1 {
		}
	}
}

func BenchmarkBigStringParse_JWT_GO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pass, err := v2.ParseRefreshToken(bigStringTokenV2, bkey1)
		if err != nil && len(pass) < 1 {
		}
	}
}

func BenchmarkBigStringParse_JWT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pass, err := v3.ParseRefreshToken(bigStringTokenV3)
		if err != nil && len(pass) < 1 {
		}
	}
}

func BenchmarkBigStringParse_FST(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token, err := converter.ParseToken(bigStringTokenFST)
		if err != nil && len(token) < 1 {
		}
	}
}

func BenchmarkUintGen_GoJose_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v5.NewAccessToken(id)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkUintGen_GolangJWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v4.NewAccessToken(id, bkey1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkUintGen_JWT_GO_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v2.NewAccessToken(id, bkey1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkUintGen_JWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v3.NewAccessToken(id)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkUintGen_FST_ASYNC(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	bid := U2B(id)
	b.ResetTimer()
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token := converter.NewToken(bid)
			if len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringGen_GoJose_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v5.NewRefreshToken(message1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringGen_GolangJWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v4.NewRefreshToken(message1, bkey1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringGen_JWT_GO_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v2.NewRefreshToken(bkey1, message1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringGen_JWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := v3.NewRefreshToken(message1)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringGen_FST_ASYNC(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	bmessage := []byte(message1)
	b.ResetTimer()
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token := converter.NewToken(bmessage)
			if len(token) < 1 {
			}
		}
	})
}

func BenchmarkUintParse_GoJose_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idT, err := v5.ParseAccessToken(uintTokenV5)
			if err != nil && idT < 1 {
			}
		}
	})
}

func BenchmarkUintParse_GolangJWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idT, err := v4.ParseAccessToken(uintTokenV4, bkey1)
			if err != nil && idT < 1 {
			}
		}
	})
}

func BenchmarkUintParse_JWT_GO_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idT, err := v2.ParseAccessToken(uintTokenV2, bkey1)
			if err != nil && idT < 1 {
			}
		}
	})
}

func BenchmarkUintParse_JWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			idT, err := v3.ParseAccessToken(uintTokenV3)
			if err != nil && idT < 1 {
			}
		}
	})
}

func BenchmarkUintParse_FST_ASYNC(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	b.ResetTimer()
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := converter.ParseToken(uintTokenFST)
			if err != nil && len(token) < 1 {
			}
		}
	})
}

func BenchmarkBigStringParse_GoJose_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pass, err := v5.ParseRefreshToken(bigStringTokenV5)
			if err != nil && len(pass) < 1 {
			}
		}
	})
}

func BenchmarkBigStringParse_GolangJWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pass, err := v4.ParseRefreshToken(bigStringTokenV4, bkey1)
			if err != nil && len(pass) < 1 {
			}
		}
	})
}

func BenchmarkBigStringParse_JWT_GO_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pass, err := v2.ParseRefreshToken(bigStringTokenV2, bkey1)
			if err != nil && len(pass) < 1 {
			}
		}
	})
}

func BenchmarkBigStringParse_JWT_ASYNC(b *testing.B) {
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pass, err := v3.ParseRefreshToken(bigStringTokenV3)
			if err != nil && len(pass) < 1 {
			}
		}
	})
}

func BenchmarkBigStringParse_FST_ASYNC(b *testing.B) {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: bkey1,
	})
	b.ResetTimer()
	b.SetParallelism(128)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			token, err := converter.ParseToken(bigStringTokenFST)
			if err != nil && len(token) < 1 {
			}
		}
	})
}
