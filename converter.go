// Package fst is a high-performance, low-memory library for generating and parsing Fast Signed Tokens (FST). FST provides an alternative to JSON-based tokens and allows you to store any information that can be represented as []byte. You can use FST for the same purposes as JWT.
package fst

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"hash"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// InvalidTokenFormat means that the token is malformed.
	InvalidTokenFormat = errors.New("Invalid token format")
	// InvalidSignature means that the token is forged.
	InvalidSignature = errors.New("Invalid signature")
	// TokenExpired means that the token is expired.
	TokenExpired = errors.New("Token expired")
)

// Converter represents a token converter that can generate and parse Fast Signed Tokens.
//
// Example:
//
//	converter := fst.NewConverter(&fst.ConverterConfig{
//		SecretKey: []byte(`secret`),
//		HashType:  sha256.New,
//	})
//
//	token := converter.NewToken([]byte(`token`))
//	fmt.Println(string(token)) // s♣�♠����▬]>¶4s\n'�a→Jtoken
//
//	value, err := converter.ParseToken(token)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(string(value)) // token
//
//	converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
//		SecretKey:          []byte(`secret`),
//		Postfix:            nil,
//		ExpirationTime:     time.Minute * 5,
//		HashType:           sha256.New,
//		WithExpirationTime: true,
//	})
//
//	tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
//	fmt.Println(string(tokenWithEx)) // Something like 3d��I�j�token4n.<� ?�↨��♣u
//
//	value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(string(value)) // token
type Converter struct {
	expirationTime      atomic.Int64
	expirationTimeBytes atomic.Value
	timeBeforeExpire    int64

	secretKey []byte
	postfix   []byte

	hmacPool sync.Pool
	hashType hash.Hash

	// NewToken creates a new FST with the provided value.
	NewToken func([]byte) []byte

	// ParseToken parses a FST and returns the value.
	//
	// It can return errors like InvalidTokenFormat, InvalidSignature, TokenExpired.
	ParseToken func([]byte) ([]byte, error)
}

// ConverterConfig represents the configuration options for creating a new Converter.
//
// SecretKey is the secret used to sign the token.
//
// Postfix is the postfix to add to the token to more secure the token.
//
// ExpirationTime is the expiration time of the token. It is -1 by default and will not expire.
//
// HashType is the hash function used to sign the token.
type ConverterConfig struct {
	// SecretKey is the secret used to sign the token.
	SecretKey []byte
	// Postfix is the postfix to add to the token to more secure the token.
	Postfix []byte
	// ExpirationTime is the expiration time of the token.
	ExpirationTime time.Duration
	// HashType is the hash function used to sign the token.
	HashType func() hash.Hash

	WithExpirationTime bool
}

// NewConverter creates a new instance of the Converter based on the provided fst.ConverterConfig.
//
// Example of the usage:
//
//		converter := fst.NewConverter(&fst.ConverterConfig{
//	     SecretKey:      []byte(`secret`),
//	     Postfix:        nil,
//	     ExpirationTime: time.Minute * 5,
//	     HashType:       sha256.New,
//	     DisableLogs:    false,
//	 })
func NewConverter(cfg *ConverterConfig) *Converter {
	if !cfg.WithExpirationTime {
		cfg.ExpirationTime = -1
	}
	if cfg.HashType == nil {
		cfg.HashType = sha256.New
	}
	converter := &Converter{
		secretKey:        cfg.SecretKey,
		postfix:          cfg.Postfix,
		timeBeforeExpire: int64(cfg.ExpirationTime.Seconds()),

		hmacPool: sync.Pool{
			New: func() interface{} {
				return hmac.New(cfg.HashType, cfg.SecretKey)
			},
		},
	}

	if cfg.ExpirationTime != -1 {
		if cfg.Postfix == nil {
			converter.NewToken = converter.newTokenWithExpire
			converter.ParseToken = converter.parseTokenWithExpire
		} else {
			converter.NewToken = converter.newTokenWithExpireAndPostfix
			converter.ParseToken = converter.parseTokenWithExpireAndPostfix
		}

		converter.expirationTime.Store(time.Now().Unix())
		converter.expirationTimeBytes.Store(getBytesForInt64(time.Now().Unix()))
		go func() {
			var now int64
			for {
				time.Sleep(1 * time.Second)
				now = time.Now().Unix()
				converter.expirationTime.Store(now)
				converter.expirationTimeBytes.Store(getBytesForInt64(now))
			}
		}()
	} else {
		if cfg.Postfix == nil {
			converter.NewToken = converter.newToken
			converter.ParseToken = converter.parseToken
		} else {
			converter.NewToken = converter.newTokenWithPostfix
			converter.ParseToken = converter.parseTokenWithPostfix
		}
	}

	return converter
}

func (c *Converter) newToken(value []byte) []byte {
	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(value)
	signature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	token := make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)))
	token = append(token, getBytesFromLen(len(signature))...)
	token = append(token, signature...)
	token = append(token, value...)

	return token
}

func (c *Converter) newTokenWithExpire(value []byte) []byte {
	exTime := c.expirationTimeBytes.Load().([]byte)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(value)
	mac.Write(exTime)
	signature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	token := make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)+8))
	token = append(token, exTime...)
	token = append(token, getBytesFromLen(len(signature))...)
	token = append(token, signature...)
	token = append(token, value...)

	return token
}

func (c *Converter) newTokenWithPostfix(value []byte) []byte {
	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(value, c.postfix...))
	signature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	token := make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)))
	token = append(token, getBytesFromLen(len(signature))...)
	token = append(token, signature...)
	token = append(token, value...)

	return token
}

func (c *Converter) newTokenWithExpireAndPostfix(value []byte) []byte {
	exTime := c.expirationTimeBytes.Load().([]byte)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(value, c.postfix...))
	mac.Write(exTime)
	signature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	token := make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)+8))
	token = append(token, exTime...)
	token = append(token, getBytesFromLen(len(signature))...)
	token = append(token, signature...)
	token = append(token, value...)

	return token
}

func (c *Converter) parseToken(token []byte) ([]byte, error) {
	if len(token) < 3 {
		return nil, InvalidTokenFormat
	}

	signatureLen, signatureSize := getLenAndSize(token)
	offset := signatureSize + signatureLen
	expectedSignature := token[signatureSize:offset]
	payload := token[offset:]

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(payload)
	actualSignature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return payload, nil
}

func (c *Converter) parseTokenWithExpire(token []byte) ([]byte, error) {
	if len(token) < 11 {
		return nil, InvalidTokenFormat
	}

	exTime := getInt64(token)
	if exTime < c.expirationTime.Load() {
		return nil, TokenExpired
	}

	signatureLen, signatureSize := getLenAndSize(token[8:])
	offset := signatureSize + signatureLen + 8
	expectedSignature := token[signatureSize+8 : offset]
	payload := token[offset:]

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(payload)
	mac.Write(token[:8])
	actualSignature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return payload, nil
}

func (c *Converter) parseTokenWithPostfix(token []byte) ([]byte, error) {
	if len(token) < 3 {
		return nil, InvalidTokenFormat
	}

	signatureLen, signatureSize := getLenAndSize(token)
	offset := signatureSize + signatureLen
	expectedSignature := token[signatureSize:offset]
	payload := token[offset:]

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(payload, c.postfix...))
	actualSignature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return payload, nil
}

func (c *Converter) parseTokenWithExpireAndPostfix(token []byte) ([]byte, error) {
	if len(token) < 11 {
		return nil, InvalidTokenFormat
	}

	exTime := getInt64(token)
	if exTime < c.expirationTime.Load() {
		return nil, TokenExpired
	}

	signatureLen, signatureSize := getLenAndSize(token[8:])
	offset := signatureSize + signatureLen + 8
	expectedSignature := token[signatureSize+8 : offset]
	payload := token[offset:]

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(payload, c.postfix...))
	mac.Write(token[:8])
	actualSignature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return payload, nil
}

// SecretKey returns the secret key used by the Converter.
func (c *Converter) SecretKey() []byte {
	return c.secretKey
}

// Postfix returns the postfix used by the Converter.
func (c *Converter) Postfix() []byte {
	return c.postfix
}

// ExpireTime returns the expiration time used by the Converter.
func (c *Converter) ExpireTime() time.Duration {
	return time.Duration(c.timeBeforeExpire)
}
