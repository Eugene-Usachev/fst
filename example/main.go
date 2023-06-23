package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/Eugene-Usachev/fst"
	"time"
)

func main() {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
		DisableLogs:    false,
	})

	token := converter.NewToken([]byte(`token`))
	fmt.Println(token)

	value, err := converter.ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value))

	tokenWithEx := converter.NewTokenWithExpire([]byte(`token`))
	fmt.Println(tokenWithEx)

	value, err = converter.ParseTokenWithExpire(tokenWithEx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value))
}
