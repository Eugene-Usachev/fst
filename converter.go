//Package fst is a high-performance, low-memory library for generating and parsing Fast Signed Tokens (FST). FST provides an alternative to JSON-based tokens and allows you to store any information that can be represented as []byte. You can use FST for the same purposes as JWT.
package fst

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
	"log"
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
// timeBeforeExpire is the lifetime of the token.
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
}

// ConverterConfig represents the configuration options for creating a new Converter.
//
// SecretKey is the secret used to sign the token.
//
// Postfix is the postfix to add to the token to more secure the token.
//
// ExpirationTime is the expiration time of the token.
//
// HashType is the hash function used to sign the token.
//
// DisableLogs is a flag to disable logs. Be better to disable logs in production but use it for development. DisableLogs = false will not slow down the program.
type ConverterConfig struct {
	// SecretKey is the secret used to sign the token.
	SecretKey []byte
	// Postfix is the postfix to add to the token to more secure the token.
	Postfix []byte
	// ExpirationTime is the expiration time of the token.
	ExpirationTime time.Duration
	// HashType is the hash function used to sign the token.
	HashType func() hash.Hash

	// DisableLogs is a flag to disable logs. Be better to disable logs in production but use it for development. DisableLogs = false will not slow down the program.
	DisableLogs bool
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
	if cfg.ExpirationTime == 0 {
		if !cfg.DisableLogs {
			log.Println(`[Warning] Passed empty ExpirationTime for Converter! Set to 5 minutes by default.`)
		}
		cfg.ExpirationTime = 5 * time.Minute
	}
	if cfg.HashType == nil {
		cfg.HashType = sha256.New
		if !cfg.DisableLogs {
			log.Println(`[Warning] Passed empty HashType for Converter! Set to sha256.New by default.`)
		}
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

	converter.timeNow.Store(time.Now().Unix())
	converter.expirationTime.Store(strconv.FormatInt(time.Now().Add(converter.timeBeforeExpire).Unix(), 10))

	go func() {
		var ex64 int64
		for {
			time.Sleep(1 * time.Second)
			converter.timeNow.Add(1)
			ex64, _ = strconv.ParseInt(converter.expirationTime.Load().(string), 10, 64)
			ex64++
			converter.expirationTime.Store(strconv.FormatInt(ex64, 10))
		}
	}()
	return converter
}

// NewToken generates a new token based on the provided value.
//
// Example of the usage:
//
// converter := fst.NewConverter(<your fst.ConverterConfig>)
//
// token, err := converter.NewToken(<value>)
func (c *Converter) NewToken(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	if c.postfix != nil {
		mac.Write(append([]byte(payloadBase64), c.postfix...))
	} else {
		mac.Write([]byte(payloadBase64))
	}
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64}, ".")
}

// NewTokenWithExpire generates a new token with expiration time based on the provided value.
//
// Example of the usage:
//
// converter := fst.NewConverter(<your fst.ConverterConfig>)
//
// token, err := converter.NewTokenWithExpire(<value>)
func (c *Converter) NewTokenWithExpire(value []byte) string {
	// Create the payload
	payloadBase64 := base64.RawURLEncoding.EncodeToString(value)

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	if c.postfix != nil {
		mac.Write(append([]byte(payloadBase64), c.postfix...))
	} else {
		mac.Write([]byte(payloadBase64))
	}
	signature := mac.Sum(nil)

	c.hmacPool.Put(mac)

	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	return strings.Join([]string{payloadBase64, signatureBase64, c.expirationTime.Load().(string)}, ".")
}

var (
	InvalidTokenFormat = errors.New("Invalid token format")
	InvalidSignature   = errors.New("Invalid signature")
	TokenExpired       = errors.New("Token expired")
)

// ParseToken parses the provided token and returns the decoded value.
//
// Warning! If you pass a token that has expiration time, this function will return error and nil value!
//
// Example of the usage:
//
// token := <some token>
//
// value, err := converter.ParseToken(token)
func (c *Converter) ParseToken(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 2 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	if c.postfix != nil {
		mac.Write(append([]byte(components[0]), c.postfix...))
	} else {
		mac.Write([]byte(components[0]))
	}

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

// ParseTokenWithExpire parses the provided token with expiration time and returns the decoded value.
//
// Warning! If you pass a token that has no expiration time, this function will return error and nil value!
//
// Example of the usage:
//
// token := <some token that has expiration time>
//
// value, err := converter.ParseTokenWithExpire(token)
func (c *Converter) ParseTokenWithExpire(token string) ([]byte, error) {
	components := strings.Split(token, ".")

	if len(components) != 3 {
		return nil, InvalidTokenFormat
	}

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()

	if c.postfix != nil {
		mac.Write(append([]byte(components[0]), c.postfix...))
	} else {
		mac.Write([]byte(components[0]))
	}

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
