package fst

import (
	"encoding/base64"
)

// EncodedConverter represents a token converter that can generate and parse Fast Signed Tokens that are encoded.
//
// # Attention!
//
// Browsers can use this, but if you don't need it, you can use Converter instead, as it is faster and more lightweight.
//
// # Example:
//
//	converter := fst.NewEncodedConverter(&fst.ConverterConfig{
//			SecretKey: []byte(`secret`),
//			HashType:  sha256.New,
//		})
//
//		token := converter.NewToken([]byte(`token`))
//		fmt.Println(token) //IOlBEQ49K_6CYh8OPhQ0cw1zBdEGxfaMhxZdCyekYRpKdG9rZW4=
//
//		value, err := converter.ParseToken(token)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(string(value)) // token
//
//		converterWithExpirationTime := fst.NewEncodedConverter(&fst.ConverterConfig{
//			SecretKey:      []byte(`secret`),
//			Postfix:        nil,
//			ExpirationTime: time.Minute * 5,
//			HashType:       sha256.New,
//		})
//
//		tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
//		fmt.Println(tokenWithEx) // Something like azriZQAAAAAgNujiWdAI6NmfGiWnt3YRNXSD1ivpdhb-8-Y8_bIIK7h0b2tlbg==
//
//		value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(string(value)) // token
type EncodedConverter struct {
	converter *Converter
}

// NewEncodedConverter creates a new instance of the EncodedConverter based on the provided fst.ConverterConfig.
//
// Example of the usage:
//
//	converter := fst.NewEncodedConverter(&fst.ConverterConfig{
//	    SecretKey:      []byte(`secret`),
//	    Postfix:        []byte(`postfix`),
//	    ExpirationTime: time.Minute * 5,
//	    HashType:       sha256.New,
//	})
func NewEncodedConverter(cfg *ConverterConfig) *EncodedConverter {
	return &EncodedConverter{
		converter: NewConverter(cfg),
	}
}

// NewToken creates a new FST with the provided value. This method encodes the token in base64.
func (c *EncodedConverter) NewToken(id []byte) string {
	return base64.URLEncoding.EncodeToString(c.converter.NewToken(id))
}

// ParseToken parses a FST and returns the value.
// This method will copy the token's value.
//
// It can return errors like InvalidTokenFormat, InvalidSignature, TokenExpired.
func (c *EncodedConverter) ParseToken(token string) ([]byte, error) {
	decodedToken, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	return c.converter.ParseToken(decodedToken)
}
