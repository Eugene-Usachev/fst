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

type Converter struct {
	timeNow          atomic.Int64
	timeExpire       atomic.Value
	timeBeforeExpire time.Duration

	secretKey []byte
	postfix   []byte

	hmacPool sync.Pool
	hashType hash.Hash
}

type ConverterConfig struct {
	SecretKey  []byte
	Postfix    []byte
	ExpireTime time.Duration
	HashType   func() hash.Hash

	DisableLogs bool
}

func NewConverter(cfg *ConverterConfig) *Converter {
	if cfg.ExpireTime == 0 {
		if !cfg.DisableLogs {
			log.Println(`[Warning] Passed empty ExpireTime for Converter! Set to 5 minutes by default.`)
		}
		cfg.ExpireTime = 5 * time.Minute
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
		timeBeforeExpire: cfg.ExpireTime,

		hmacPool: sync.Pool{
			New: func() interface{} {
				return hmac.New(cfg.HashType, cfg.SecretKey)
			},
		},
	}

	converter.timeNow.Store(time.Now().Unix())
	converter.timeExpire.Store(strconv.FormatInt(time.Now().Add(converter.timeBeforeExpire).Unix(), 10))

	go func() {
		var ex64 int64
		for {
			time.Sleep(1 * time.Second)
			converter.timeNow.Add(1)
			ex64, _ = strconv.ParseInt(converter.timeExpire.Load().(string), 10, 64)
			ex64++
			converter.timeExpire.Store(strconv.FormatInt(ex64, 10))
		}
	}()
	return converter
}

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

	return strings.Join([]string{payloadBase64, signatureBase64, c.timeExpire.Load().(string)}, ".")
}

var (
	InvalidTokenFormat = errors.New("Invalid token format")
	InvalidSignature   = errors.New("Invalid signature")
	TokenExpired       = errors.New("Token expired")
)

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

func (c *Converter) SecretKey() []byte {
	return c.secretKey
}

func (c *Converter) Postfix() []byte {
	return c.postfix
}

func (c *Converter) ExpireTime() time.Duration {
	return c.timeBeforeExpire
}
