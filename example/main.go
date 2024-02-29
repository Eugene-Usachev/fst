package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/Eugene-Usachev/fst"
	"time"
)

func main() {
	converter := fst.NewConverter(&fst.ConverterConfig{
		SecretKey: []byte(`secret`),
		HashType:  sha256.New,
	})

	token := converter.NewToken([]byte(`token`))
	fmt.Println(string(token)) //+XLyVq4BwH3MMLunFEboXU2OvtYAEgnlGsvwWXibAdG0dG9rZW4

	value, err := converter.ParseToken(token)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token

	converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
		SecretKey:          []byte(`secret`),
		Postfix:            nil,
		ExpirationTime:     time.Minute * 5,
		HashType:           sha256.New,
		WithExpirationTime: true,
	})

	tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
	fmt.Println(string(tokenWithEx)) // ˼�e+il7xQKjrk9p3CwFhDlqQziBfuqrEdgoKK-hKFmyCU0IdG9rZW4

	value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(value)) // token
}
