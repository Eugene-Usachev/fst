// Package fst is a high-performance, low-memory library for generating and parsing Fast Signed Tokens (FST). FST provides an alternative to JSON-based tokens and allows you to store any information that can be represented as []byte. You can use FST for the same purposes as JWT.
package fst

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/Eugene-Usachev/fastbytes"
	"hash"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Converter represents a token converter that can generate and parse Fast Signed Tokens.
//
// secretKey is the secret used to sign the token.
//
// postfix is the postfix to add to the token to more secure the token.
//
// hashType is the hash function used to sign the token.
//
// timeBeforeExpire is the lifetime of the token. By default, it is -1
//
// hmacPool and expirationTime and timeNow are needed to improve performance.
type Converter struct {
	timeNow          atomic.Int64
	expirationTime   atomic.Value
	timeBeforeExpire time.Duration

	secretKey []byte
	postfix   []byte

	hmacPool sync.Pool
	hashType hash.Hash

	NewToken   func([]byte) string
	ParseToken func(string) ([]byte, error)
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
		timeBeforeExpire: cfg.ExpirationTime,

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

		converter.timeNow.Store(time.Now().Unix())
		converter.expirationTime.Store(strconv.FormatInt(time.Now().Add(converter.timeBeforeExpire).Unix(), 10))

		go func() {
			for {
				time.Sleep(1 * time.Second)
				converter.timeNow.Store(time.Now().Unix())
				converter.expirationTime.Store(strconv.FormatInt(time.Now().Add(converter.timeBeforeExpire).Unix(), 10))
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

func (c *Converter) newToken(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(fastbytes.S2B(payloadBase64))
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64}, ".")
}

func (c *Converter) newTokenWithExpire(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)
	exTime := c.expirationTime.Load().(string)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(fastbytes.S2B(payloadBase64 + exTime))
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64, exTime}, ".")
}

func (c *Converter) newTokenWithPostfix(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(fastbytes.S2B(payloadBase64), c.postfix...))
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64}, ".")
}

func (c *Converter) newTokenWithExpireAndPostfix(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)
	exTime := c.expirationTime.Load().(string)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(fastbytes.S2B(payloadBase64+exTime), c.postfix...))
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64, exTime}, ".")
}

var (
	InvalidTokenFormat = errors.New("Invalid token format")
	InvalidSignature   = errors.New("Invalid signature")
	TokenExpired       = errors.New("Token expired")
)

func (c *Converter) parseToken(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 2 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(fastbytes.S2B(components[0]))

	expectedSignature, err := base64.RawURLEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	actualSignature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return base64.RawURLEncoding.DecodeString(components[0])
}

func (c *Converter) parseTokenWithExpire(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 3 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(fastbytes.S2B(components[0] + components[2]))

	expectedSignature, err := base64.RawURLEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	actualSignature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	expiration, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return nil, err
	}

	if c.timeNow.Load() > expiration {
		return nil, TokenExpired
	}

	return base64.RawURLEncoding.DecodeString(components[0])
}

func (c *Converter) parseTokenWithPostfix(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 2 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(fastbytes.S2B(components[0]), c.postfix...))

	expectedSignature, err := base64.RawURLEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	actualSignature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	return base64.RawURLEncoding.DecodeString(components[0])
}

func (c *Converter) parseTokenWithExpireAndPostfix(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 3 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(append(fastbytes.S2B(components[0]+components[2]), c.postfix...))

	expectedSignature, err := base64.RawURLEncoding.DecodeString(components[1])
	if err != nil {
		return nil, err
	}

	actualSignature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	if !hmac.Equal(expectedSignature, actualSignature) {
		return nil, InvalidSignature
	}

	expiration, err := strconv.ParseInt(components[2], 10, 64)
	if err != nil {
		return nil, err
	}

	if c.timeNow.Load() > expiration {
		return nil, TokenExpired
	}

	return base64.RawURLEncoding.DecodeString(components[0])
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
	return c.timeBeforeExpire
}
