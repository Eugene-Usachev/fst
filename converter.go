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
// # Be careful!
//
// Browsers cannot use this!
// To work with the browser and HTTP in general, use EncodedConverter!
//
// # Example:
//
//	converter := fst.NewConverter(&fst.ConverterConfig{
//			SecretKey: []byte(`secret`),
//			HashType:  sha256.New,
//		})
//
//		token := converter.NewToken([]byte(`token`))
//		fmt.Println(string(token)) //s♣�♠����▬]>¶4s\n'�a→Jtoken
//
//		value, err := converter.ParseToken(token)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(string(value)) // token
//
//		converterWithExpirationTime := fst.NewConverter(&fst.ConverterConfig{
//			SecretKey:      []byte(`secret`),
//			Postfix:        nil,
//			ExpirationTime: time.Minute * 5,
//			HashType:       sha256.New,
//		})
//
//		tokenWithEx := converterWithExpirationTime.NewToken([]byte(`token`))
//		fmt.Println(string(tokenWithEx)) // Something like k:�e 6��Y�ٟ→%��v◄5t��+�v▬���<�+�token
//
//		value, err = converterWithExpirationTime.ParseToken(tokenWithEx)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(string(value)) // token
type Converter struct {
	expirationTime      atomic.Int64
	expirationTimeBytes atomic.Value
	timeBeforeExpire    int64

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
// ExpirationTime is the expiration time of the token. It is zero by default and will not expire.
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
}

// NewConverter creates a new instance of the Converter based on the provided fst.ConverterConfig.
//
// Example of the usage:
//
//	converter := fst.NewConverter(&fst.ConverterConfig{
//	    SecretKey:      []byte(`secret`),
//	    Postfix:        []byte(`postfix`),
//	    ExpirationTime: time.Minute * 5,
//	    HashType:       sha256.New,
//	})
func NewConverter(cfg *ConverterConfig) *Converter {
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

	if cfg.ExpirationTime != 0 {
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
	}

	return converter
}

// NewToken creates a new FST with the provided value. This method does not encode the token in base64.
func (c *Converter) NewToken(value []byte) []byte {
	var exTime []byte
	isWithExpirationTime := c.timeBeforeExpire != 0

	// Create the signature
	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(value)
	if isWithExpirationTime {
		exTime = c.expirationTimeBytes.Load().([]byte)
		mac.Write(exTime)
	}
	if c.postfix != nil {
		mac.Write(c.postfix)
	}
	signature := mac.Sum(nil)
	c.hmacPool.Put(mac)

	var token []byte

	if !isWithExpirationTime {
		token = make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)))
		token = append(token, getBytesFromLen(len(signature))...)
		token = append(token, signature...)
		token = append(token, value...)
	} else {
		token = make([]byte, 0, len(value)+len(signature)+getSizeForLen(len(signature)+8))
		token = append(token, exTime...)
		token = append(token, getBytesFromLen(len(signature))...)
		token = append(token, signature...)
		token = append(token, value...)
	}

	return token
}

// ParseToken parses a FST and returns the value.
// This method will use token to return the value, instead of copying.
//
// It can return errors like InvalidTokenFormat, InvalidSignature, TokenExpired.
func (c *Converter) ParseToken(token []byte) ([]byte, error) {
	if len(token) < 11 && (c.timeBeforeExpire == 0 && len(token) < 3) {
		panic(len(token))
		return nil, InvalidTokenFormat
	}

	isWithExpirationTime := c.timeBeforeExpire != 0

	var payloadOffset int
	if isWithExpirationTime {
		exTime := getInt64(token)
		if exTime < c.expirationTime.Load() {
			return nil, TokenExpired
		}
		payloadOffset = 8
	}
	signatureLen, signatureSize := getLenAndSize(token[payloadOffset:])
	signatureOffset := payloadOffset + signatureSize
	payloadOffset += signatureSize + signatureLen

	if len(token) <= payloadOffset {
		return nil, InvalidTokenFormat
	}

	expectedSignature := token[signatureOffset:payloadOffset]
	payload := token[payloadOffset:]

	mac := c.hmacPool.Get().(hash.Hash)
	mac.Reset()
	mac.Write(payload)
	if isWithExpirationTime {
		mac.Write(token[:8])
	}
	if c.postfix != nil {
		mac.Write(c.postfix)
	}
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
