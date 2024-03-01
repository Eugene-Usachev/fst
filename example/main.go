package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/Eugene-Usachev/fst"
	"time"
)

func main() {
	rawConvert()
	encodedConvert()
}

func rawConvert() {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: []byte(`secret`),
		HashType:  sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	fmt.Println(string(token)) //s♣�♠����▬]>¶4s\n'�a→Jtoken

	value, err := converter.ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token

	converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})

	tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
	fmt.Println(string(tokenWithEx)) // Something like k:�e 6��Y�ٟ→%��v◄5t��+�v▬���<�+�token

	value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token
}

func encodedConvert() {
	converter := fst.NewEncodedConverter(&fst.ConverterConfig{
		SecretKey: []byte(`secret`),
		HashType:  sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	fmt.Println(token) //IOlBEQ49K_6CYh8OPhQ0cw1zBdEGxfaMhxZdCyekYRpKdG9rZW4=

	value, err := converter.ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token

	converterWithExpirationTime := fst.NewEncodedConverter(&fst.ConverterConfig{
		SecretKey:      []byte(`secret`),
		Postfix:        nil,
		ExpirationTime: time.Minute * 5,
		HashType:       sha256.New,
	})

	tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
	fmt.Println(tokenWithEx) // Something like azriZQAAAAAgNujiWdAI6NmfGiWnt3YRNXSD1ivpdhb-8-Y8_bIIK7h0b2tlbg==

	value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token
}
